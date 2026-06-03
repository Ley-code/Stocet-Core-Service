package controllers

import (
	"stocet-core-service/internal/database"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type AdminAnalyticsController struct{}

func NewAdminAnalyticsController() *AdminAnalyticsController {
	return &AdminAnalyticsController{}
}

type userTierRow struct {
	Tier  string `gorm:"column:tier"`
	Count int64  `gorm:"column:count"`
}

type aiDayRow struct {
	Day      string `gorm:"column:day"`
	Requests int64  `gorm:"column:requests"`
	Users    int64  `gorm:"column:users"`
}

// Get godoc — GET /admin/analytics
func (c *AdminAnalyticsController) Get(ctx *gin.Context) {
	db := database.DB

	// ── User counts ──────────────────────────────────────────────────────────
	var totalUsers, verifiedUsers, newUsers7d, newUsers30d int64
	db.Raw(`SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`).Scan(&totalUsers)
	db.Raw(`SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND is_verified = true`).Scan(&verifiedUsers)
	db.Raw(`SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND created_at >= NOW() - INTERVAL '7 days'`).Scan(&newUsers7d)
	db.Raw(`SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND created_at >= NOW() - INTERVAL '30 days'`).Scan(&newUsers30d)

	// Users by subscription tier — initialise to empty slice so JSON is [] not null
	tierRows := make([]userTierRow, 0)
	db.Raw(`
		SELECT subscription_status AS tier, COUNT(*) AS count
		FROM users
		WHERE deleted_at IS NULL
		GROUP BY subscription_status
	`).Scan(&tierRows)

	// ── AI quota stats ───────────────────────────────────────────────────────
	var aiToday, ai7d int64
	db.Raw(`SELECT COALESCE(SUM(request_count), 0) FROM ai_usage_daily WHERE usage_date = CURRENT_DATE`).Scan(&aiToday)
	db.Raw(`SELECT COALESCE(SUM(request_count), 0) FROM ai_usage_daily WHERE usage_date >= CURRENT_DATE - 6`).Scan(&ai7d)

	// Active AI users today
	var aiUsersToday int64
	db.Raw(`SELECT COUNT(DISTINCT user_id) FROM ai_usage_daily WHERE usage_date = CURRENT_DATE AND request_count > 0`).Scan(&aiUsersToday)

	// Daily breakdown for last 7 days
	aiByDay := make([]aiDayRow, 0)
	db.Raw(`
		SELECT usage_date::text AS day,
		       SUM(request_count) AS requests,
		       COUNT(DISTINCT user_id) AS users
		FROM ai_usage_daily
		WHERE usage_date >= CURRENT_DATE - 6
		GROUP BY usage_date
		ORDER BY usage_date ASC
	`).Scan(&aiByDay)

	// ── Feature engagement ───────────────────────────────────────────────────
	var portfolioUsers, totalHoldings int64
	db.Raw(`SELECT COUNT(DISTINCT user_id) FROM portfolio_holdings`).Scan(&portfolioUsers)
	db.Raw(`SELECT COUNT(*) FROM portfolio_holdings`).Scan(&totalHoldings)

	var totalAlerts, activeAlerts int64
	db.Raw(`SELECT COUNT(*) FROM alerts WHERE deleted_at IS NULL`).Scan(&totalAlerts)
	db.Raw(`SELECT COUNT(*) FROM alerts WHERE deleted_at IS NULL AND status = 'active'`).Scan(&activeAlerts)

	var watchlistUsers, totalWatchlistItems int64
	db.Raw(`SELECT COUNT(DISTINCT user_id) FROM watchlists`).Scan(&watchlistUsers)
	db.Raw(`SELECT COUNT(*) FROM watchlists`).Scan(&totalWatchlistItems)

	var advisoryReports int64
	db.Raw(`SELECT COUNT(*) FROM advisory_reports`).Scan(&advisoryReports)

	// ── Calendar ─────────────────────────────────────────────────────────────
	var totalCorpActions, upcomingCorpActions int64
	db.Raw(`SELECT COUNT(*) FROM corporate_actions`).Scan(&totalCorpActions)
	db.Raw(`SELECT COUNT(*) FROM corporate_actions WHERE event_date BETWEEN CURRENT_DATE AND CURRENT_DATE + 30`).Scan(&upcomingCorpActions)

	var totalEarnings, upcomingEarnings int64
	db.Raw(`SELECT COUNT(*) FROM earnings_dates`).Scan(&totalEarnings)
	db.Raw(`SELECT COUNT(*) FROM earnings_dates WHERE expected_date BETWEEN CURRENT_DATE AND CURRENT_DATE + 30`).Scan(&upcomingEarnings)

	// ── DAU / MAU ─────────────────────────────────────────────────────────────
	// Proxy: union of users active via AI usage or alert triggers (best available signal).
	var dau, mau int64
	db.Raw(`
		SELECT COUNT(*) FROM (
			SELECT user_id FROM ai_usage_daily
			WHERE usage_date = CURRENT_DATE AND request_count > 0
			UNION
			SELECT user_id FROM alert_trigger_history
			WHERE triggered_at >= NOW() - INTERVAL '1 day'
		) active_today
	`).Scan(&dau)
	db.Raw(`
		SELECT COUNT(*) FROM (
			SELECT user_id FROM ai_usage_daily
			WHERE usage_date >= CURRENT_DATE - 29 AND request_count > 0
			UNION
			SELECT user_id FROM alert_trigger_history
			WHERE triggered_at >= NOW() - INTERVAL '30 days'
		) active_month
	`).Scan(&mau)

	// ── Alert fires ───────────────────────────────────────────────────────────
	var alertFiresToday, alertFires7d int64
	db.Raw(`SELECT COUNT(*) FROM alert_trigger_history WHERE triggered_at >= NOW() - INTERVAL '1 day'`).Scan(&alertFiresToday)
	db.Raw(`SELECT COUNT(*) FROM alert_trigger_history WHERE triggered_at >= NOW() - INTERVAL '7 days'`).Scan(&alertFires7d)

	// ── Org & invite ─────────────────────────────────────────────────────────
	var totalOrgs, totalInviteCodes, usedInviteCodes int64
	db.Raw(`SELECT COUNT(*) FROM organizations`).Scan(&totalOrgs)
	db.Raw(`SELECT COUNT(*) FROM invite_codes WHERE is_active = true`).Scan(&totalInviteCodes)
	db.Raw(`SELECT COALESCE(SUM(current_uses), 0) FROM invite_codes`).Scan(&usedInviteCodes)

	utils.SuccessResponse(ctx, 200, gin.H{
		"users": gin.H{
			"total":    totalUsers,
			"verified": verifiedUsers,
			"new_7d":   newUsers7d,
			"new_30d":  newUsers30d,
			"by_tier":  tierRows,
			"dau":      dau,
			"mau":      mau,
		},
		"ai": gin.H{
			"requests_today":     aiToday,
			"requests_7d":        ai7d,
			"active_users_today": aiUsersToday,
			"by_day":             aiByDay,
		},
		"features": gin.H{
			"portfolio_users":    portfolioUsers,
			"total_holdings":     totalHoldings,
			"total_alerts":       totalAlerts,
			"active_alerts":      activeAlerts,
			"alert_fires_today":  alertFiresToday,
			"alert_fires_7d":     alertFires7d,
			"watchlist_users":    watchlistUsers,
			"watchlist_items":    totalWatchlistItems,
			"advisory_reports":   advisoryReports,
		},
		"calendar": gin.H{
			"total_actions":    totalCorpActions,
			"upcoming_30d":     upcomingCorpActions,
			"total_earnings":   totalEarnings,
			"upcoming_earnings_30d": upcomingEarnings,
		},
		"orgs": gin.H{
			"total_orgs":        totalOrgs,
			"active_inv_codes":  totalInviteCodes,
			"total_invite_uses": usedInviteCodes,
		},
	}, "")
}
