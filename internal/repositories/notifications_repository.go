package repositories

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"
)

type NotificationsRepository struct{}

func NewNotificationsRepository() *NotificationsRepository {
	return &NotificationsRepository{}
}

// MatchAndTriggerAlerts finds active alerts whose conditions are met by the
// latest market prices, creates a notification for each match, and marks
// the alert as triggered so it does not fire again.
func (r *NotificationsRepository) MatchAndTriggerAlerts() (int, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}

	// Step 1: Find all matched alerts using the latest price per ticker.
	type matchedAlert struct {
		AlertID   string  `gorm:"column:alert_id"`
		UserID    string  `gorm:"column:user_id"`
		Ticker    string  `gorm:"column:ticker"`
		Condition string  `gorm:"column:condition"`
		Threshold float64 `gorm:"column:threshold"`
		Price     float64 `gorm:"column:price"`
	}

	var matches []matchedAlert
	err := tx.Raw(`
		SELECT a.id AS alert_id,
		       a.user_id,
		       a.ticker,
		       a.condition,
		       a.threshold,
		       lp.price
		FROM alerts a
		JOIN LATERAL (
			SELECT mp.price
			FROM market_prices mp
			WHERE mp.ticker_symbol = a.ticker
			ORDER BY mp.timestamp DESC
			LIMIT 1
		) lp ON true
		WHERE a.status = 'active'
		  AND a.deleted_at IS NULL
		  AND (
		        (a.condition = 'above' AND lp.price >= a.threshold)
		     OR (a.condition = 'below' AND lp.price <= a.threshold)
		  )
	`).Scan(&matches).Error
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to query matched alerts: %w", err)
	}

	if len(matches) == 0 {
		tx.Commit()
		return 0, nil
	}

	// Step 2: Create notifications for each matched alert.
	now := time.Now()
	for _, m := range matches {
		title := fmt.Sprintf("Price Alert: %s", m.Ticker)
		message := fmt.Sprintf(
			"%s has reached %.2f ETB (your target: %s %.2f ETB)",
			m.Ticker, m.Price, m.Condition, m.Threshold,
		)

		userUUID, err := uuid.Parse(m.UserID)
		if err != nil {
			log.Printf("Invalid user UUID %s for alert %s: %v", m.UserID, m.AlertID, err)
			continue
		}

		notification := models.Notification{
			UserID:    userUUID,
			Title:     title,
			Message:   message,
			CreatedAt: now,
		}

		if err := tx.Create(&notification).Error; err != nil {
			log.Printf("Failed to create notification for alert %s: %v", m.AlertID, err)
			continue
		}
	}

	// Step 3: Mark matched alerts as triggered.
	alertIDs := make([]string, len(matches))
	for i, m := range matches {
		alertIDs[i] = m.AlertID
	}

	if err := tx.Exec(
		`UPDATE alerts SET status = 'triggered', triggered_at = ? WHERE id IN ?`,
		now, alertIDs,
	).Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to update triggered alerts: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return len(matches), nil
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
