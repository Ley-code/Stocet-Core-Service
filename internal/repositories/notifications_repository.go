package repositories

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"stocet-core-service/internal/database"
	"stocet-core-service/internal/email"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/push"
)

type NotificationsRepository struct {
	emailSender email.Sender
	fcm         *push.FCMService
}

func NewNotificationsRepository(emailSender email.Sender, fcm *push.FCMService) *NotificationsRepository {
	return &NotificationsRepository{emailSender: emailSender, fcm: fcm}
}

// MatchAndTriggerAlerts finds all active alerts whose conditions are met by
// the latest market prices or new corporate news, creates notifications,
// records trigger history, and marks/resets the alert status.
func (r *NotificationsRepository) MatchAndTriggerAlerts() (int, error) {
	now := time.Now()
	var triggered int

	// ── 1. Price alerts (above / below) ──────────────────────────────────────
	type priceMatch struct {
		AlertID       string  `gorm:"column:alert_id"`
		UserID        string  `gorm:"column:user_id"`
		UserEmail     string  `gorm:"column:user_email"`
		Ticker        string  `gorm:"column:ticker"`
		AlertType     string  `gorm:"column:alert_type"`
		Condition     string  `gorm:"column:condition"`
		Threshold     float64 `gorm:"column:threshold"`
		Price         float64 `gorm:"column:price"`
		PercentChange float64 `gorm:"column:percent_change"`
		IsRecurring   bool    `gorm:"column:is_recurring"`
		CooldownHours int     `gorm:"column:cooldown_hours"`
	}

	var priceMatches []priceMatch
	err := database.DB.Raw(`
		SELECT a.id          AS alert_id,
		       a.user_id,
		       u.email       AS user_email,
		       a.ticker,
		       a.alert_type,
		       a.condition,
		       a.threshold,
		       lp.price,
		       lp.percent_change,
		       a.is_recurring,
		       a.cooldown_hours
		FROM alerts a
		JOIN users u ON u.id = a.user_id
		JOIN LATERAL (
			SELECT mp.price, mp.percent_change
			FROM market_prices mp
			WHERE mp.ticker_symbol = a.ticker
			ORDER BY mp.timestamp DESC
			LIMIT 1
		) lp ON true
		WHERE a.status = 'active'
		  AND a.deleted_at IS NULL
		  AND a.alert_type IN ('price', 'percent_change')
		  AND (
		    a.last_triggered_at IS NULL
		    OR EXTRACT(EPOCH FROM (NOW() - a.last_triggered_at)) / 3600 >= a.cooldown_hours
		  )
		  AND (
		    (a.alert_type = 'price'          AND a.condition = 'above' AND lp.price          >= a.threshold)
		 OR (a.alert_type = 'price'          AND a.condition = 'below' AND lp.price          <= a.threshold)
		 OR (a.alert_type = 'percent_change' AND a.condition = 'above' AND lp.percent_change >= a.threshold)
		 OR (a.alert_type = 'percent_change' AND a.condition = 'below' AND lp.percent_change <= a.threshold)
		  )
	`).Scan(&priceMatches).Error
	if err != nil {
		return 0, fmt.Errorf("price alert query: %w", err)
	}

	for _, m := range priceMatches {
		if err := r.fireAlert(fireParams{
			alertID:       m.AlertID,
			userID:        m.UserID,
			userEmail:     m.UserEmail,
			ticker:        m.Ticker,
			alertType:     m.AlertType,
			condition:     m.Condition,
			threshold:     m.Threshold,
			price:         m.Price,
			percentChange: m.PercentChange,
			isRecurring:   m.IsRecurring,
			now:           now,
		}); err != nil {
			log.Printf("[alerts] fire error alert=%s: %v", m.AlertID, err)
			continue
		}
		triggered++
	}

	// ── 2. New-news alerts ────────────────────────────────────────────────────
	type newsMatch struct {
		AlertID       string  `gorm:"column:alert_id"`
		UserID        string  `gorm:"column:user_id"`
		UserEmail     string  `gorm:"column:user_email"`
		Ticker        string  `gorm:"column:ticker"`
		NewsTitle     string  `gorm:"column:news_title"`
		Price         float64 `gorm:"column:price"`
		PercentChange float64 `gorm:"column:percent_change"`
		IsRecurring   bool    `gorm:"column:is_recurring"`
		CooldownHours int     `gorm:"column:cooldown_hours"`
	}
	var newsMatches []newsMatch
	err = database.DB.Raw(`
		SELECT a.id     AS alert_id,
		       a.user_id,
		       u.email  AS user_email,
		       a.ticker,
		       ln.title AS news_title,
		       COALESCE(lp.price, 0)          AS price,
		       COALESCE(lp.percent_change, 0) AS percent_change,
		       a.is_recurring,
		       a.cooldown_hours
		FROM alerts a
		JOIN users u ON u.id = a.user_id
		JOIN company_profiles cp ON cp.ticker_symbol = a.ticker
		JOIN LATERAL (
			SELECT cn.title
			FROM corporate_news cn
			WHERE cn.company_profile_id = cp.id
			  AND cn.published_at > COALESCE(a.last_triggered_at, a.created_at)
			ORDER BY cn.published_at DESC
			LIMIT 1
		) ln ON true
		LEFT JOIN LATERAL (
			SELECT mp.price, mp.percent_change
			FROM market_prices mp
			WHERE mp.ticker_symbol = a.ticker
			ORDER BY mp.timestamp DESC
			LIMIT 1
		) lp ON true
		WHERE a.status = 'active'
		  AND a.deleted_at IS NULL
		  AND a.alert_type = 'new_news'
		  AND (
		    a.last_triggered_at IS NULL
		    OR EXTRACT(EPOCH FROM (NOW() - a.last_triggered_at)) / 3600 >= a.cooldown_hours
		  )
	`).Scan(&newsMatches).Error
	if err != nil {
		log.Printf("[alerts] new_news query error: %v", err)
	} else {
		for _, m := range newsMatches {
			if err := r.fireAlert(fireParams{
				alertID:       m.AlertID,
				userID:        m.UserID,
				userEmail:     m.UserEmail,
				ticker:        m.Ticker,
				alertType:     "new_news",
				condition:     "",
				threshold:     0,
				price:         m.Price,
				percentChange: m.PercentChange,
				isRecurring:   m.IsRecurring,
				conditionText: fmt.Sprintf("New announcement: %s", truncateTitle(m.NewsTitle, 80)),
				now:           now,
			}); err != nil {
				log.Printf("[alerts] fire news error alert=%s: %v", m.AlertID, err)
				continue
			}
			triggered++
		}
	}

	return triggered, nil
}

type fireParams struct {
	alertID       string
	userID        string
	userEmail     string
	ticker        string
	alertType     string
	condition     string
	threshold     float64
	price         float64
	percentChange float64
	isRecurring   bool
	conditionText string // override; computed if empty
	now           time.Time
}

func (r *NotificationsRepository) fireAlert(p fireParams) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	userUUID, err := uuid.Parse(p.userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("invalid user uuid: %w", err)
	}

	// Build human-readable condition text if not provided
	condText := p.conditionText
	if condText == "" {
		condText = buildConditionText(p.alertType, p.condition, p.threshold, p.price, p.percentChange)
	}

	// Notification
	notification := models.Notification{
		UserID:  userUUID,
		Title:   fmt.Sprintf("Alert: %s", p.ticker),
		Message: condText,
	}
	if err := tx.Create(&notification).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create notification: %w", err)
	}

	// Trigger history
	history := models.AlertTriggerHistory{
		AlertID:        p.alertID,
		UserID:         userUUID,
		Ticker:         p.ticker,
		AlertType:      models.AlertType(p.alertType),
		Condition:      p.condition,
		Threshold:      p.threshold,
		TriggeredPrice: p.price,
		PercentChange:  p.percentChange,
		TriggeredAt:    p.now,
	}
	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create history: %w", err)
	}

	// Update alert status
	updates := map[string]interface{}{
		"last_triggered_at": p.now,
		"triggered_at":      p.now,
	}
	if !p.isRecurring {
		updates["status"] = "triggered"
	}
	if err := tx.Model(&models.Alert{}).Where("id = ?", p.alertID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("update alert: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	// Fetch FCM token for this user (best-effort, outside transaction)
	var fcmToken string
	if r.fcm != nil && r.fcm.Enabled() {
		var u models.User
		if err := database.DB.Select("fcm_token").Where("id = ?", p.userID).First(&u).Error; err == nil {
			if u.FCMToken != nil {
				fcmToken = *u.FCMToken
			}
		}
	}

	// Send push notification (best-effort)
	if fcmToken != "" {
		pushCtx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		pushData := map[string]string{
			"type":   "alert",
			"ticker": p.ticker,
		}
		if err := r.fcm.SendPush(pushCtx, fcmToken, fmt.Sprintf("Alert: %s", p.ticker), condText, pushData); err != nil {
			logrus.Printf("[alerts] push failed alert=%s: %v", p.alertID, err)
		}
	}

	// Send email (best-effort)
	if r.emailSender != nil && r.emailSender.Enabled() && strings.TrimSpace(p.userEmail) != "" {
		deepLink := fmt.Sprintf("https://stocet.com/markets?symbol=%s", p.ticker)
		data := email.AlertEmailData{
			Ticker:        p.ticker,
			AlertType:     p.alertType,
			ConditionText: condText,
			Price:         p.price,
			PercentChange: p.percentChange,
			TriggeredAt:   p.now,
			DeepLink:      deepLink,
		}
		subject, htmlBody := email.BuildAlertEmail(data)
		if err := r.emailSender.SendHTML(p.userEmail, subject, htmlBody); err != nil {
			logrus.Printf("[alerts] email send failed alert=%s to=%s: %v", p.alertID, p.userEmail, err)
		}
	}

	return nil
}

func buildConditionText(alertType, condition string, threshold, price, pct float64) string {
	switch alertType {
	case "price":
		dir := "above"
		if condition == "below" {
			dir = "below"
		}
		return fmt.Sprintf("Price went %s ETB %.2f (current: ETB %.2f)", dir, threshold, price)
	case "percent_change":
		dir := "up"
		if condition == "below" {
			dir = "down"
		}
		return fmt.Sprintf("Daily change went %s %.2f%% (current: %.2f%%)", dir, threshold, pct)
	case "new_news":
		return fmt.Sprintf("New corporate announcement (price: ETB %.2f)", price)
	default:
		return fmt.Sprintf("Alert triggered at ETB %.2f", price)
	}
}

func truncateTitle(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func (r *NotificationsRepository) GetUnreadByUserID(userID string) ([]models.Notification, error) {
	var notifications []models.Notification
	err := database.DB.
		Where("user_id = ? AND is_read = false", userID).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

func (r *NotificationsRepository) GetAllByUserID(userID string) ([]models.Notification, error) {
	var notifications []models.Notification
	err := database.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

func (r *NotificationsRepository) MarkRead(id string, userID string) error {
	result := database.DB.
		Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}
	return nil
}
