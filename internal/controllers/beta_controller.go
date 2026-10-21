package controllers

import (
	"net/http"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"stocet-core-service/internal/database"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// BetaController handles the invite-only beta funnel: public waitlist, in-app feedback, admin review, and a data-freshness health check.
type BetaController struct{}

func NewBetaController() *BetaController { return &BetaController{} }

func clip(s string, max int) string {
	s = strings.TrimSpace(s)
	if utf8.RuneCountInString(s) > max {
		return string([]rune(s)[:max])
	}
	return s
}

type waitlistRequest struct {
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	Role         string `json:"role"`
	Organization string `json:"organization"`
	Source       string `json:"source"`
	Note         string `json:"note"`
	// honeypot: real users never fill this hidden field
	Website string `json:"website"`
}

// JoinWaitlist godoc — POST /waitlist (public, rate limited)
func (c *BetaController) JoinWaitlist(ctx *gin.Context) {
	var req waitlistRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}
	if req.Website != "" { // bot: pretend it worked
		utils.SuccessResponse(ctx, 200, gin.H{"joined": true}, "You're on the list")
		return
	}
	addr, err := mail.ParseAddress(strings.TrimSpace(req.Email))
	if err != nil || len(addr.Address) > 254 {
		utils.BadRequest(ctx, "Please enter a valid email address")
		return
	}
	err = database.DB.Exec(`
		INSERT INTO waitlist_signups (email, full_name, role, organization, source, note)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT (lower(email)) DO NOTHING`,
		strings.ToLower(addr.Address), clip(req.FullName, 120), clip(req.Role, 80), clip(req.Organization, 120), clip(req.Source, 80), clip(req.Note, 600)).Error
	if err != nil {
		utils.InternalServerError(ctx, err, "Could not save your request")
		return
	}
	// same answer whether or not the email was already there (no account/waitlist enumeration)
	utils.SuccessResponse(ctx, 200, gin.H{"joined": true}, "You're on the list")
}

type feedbackRequest struct {
	Category string `json:"category"`
	Message  string `json:"message"`
	PageURL  string `json:"page_url"`
	Email    string `json:"email"`
}

// SubmitFeedback godoc — POST /feedback (optional auth, rate limited)
func (c *BetaController) SubmitFeedback(ctx *gin.Context) {
	var req feedbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}
	msg := clip(req.Message, 4000)
	if utf8.RuneCountInString(msg) < 5 {
		utils.BadRequest(ctx, "Please describe the issue in a few words")
		return
	}
	category := req.Category
	switch category {
	case "bug", "data_error", "idea", "other":
	default:
		category = "other"
	}
	email := ""
	if a, err := mail.ParseAddress(strings.TrimSpace(req.Email)); err == nil {
		email = strings.ToLower(a.Address)
	}
	var userID interface{}
	if id, ok := utils.GetUserID(ctx); ok {
		userID = id
	}
	err := database.DB.Exec(`
		INSERT INTO feedback (user_id, email, category, message, page_url, user_agent)
		VALUES (?, ?, ?, ?, ?, ?)`,
		userID, email, category, msg, clip(req.PageURL, 500), clip(ctx.Request.UserAgent(), 300)).Error
	if err != nil {
		utils.InternalServerError(ctx, err, "Could not save your feedback")
		return
	}
	utils.SuccessResponse(ctx, 201, gin.H{"received": true}, "Thanks. We read every message.")
}

// AdminWaitlist godoc — GET /admin/beta/waitlist
func (c *BetaController) AdminWaitlist(ctx *gin.Context) {
	rows := []map[string]interface{}{}
	if err := database.DB.Raw(`SELECT id, email, full_name, role, organization, source, note, invite_sent_at, created_at FROM waitlist_signups ORDER BY created_at DESC LIMIT 500`).Scan(&rows).Error; err != nil {
		utils.InternalServerError(ctx, err, "Failed to load waitlist")
		return
	}
	utils.SuccessResponse(ctx, 200, gin.H{"signups": rows, "total": len(rows)}, "")
}

// AdminMarkInvited godoc — PATCH /admin/beta/waitlist/:id/invited
func (c *BetaController) AdminMarkInvited(ctx *gin.Context) {
	res := database.DB.Exec(`UPDATE waitlist_signups SET invite_sent_at = now() WHERE id = ?`, ctx.Param("id"))
	if res.Error != nil {
		utils.InternalServerError(ctx, res.Error, "Failed to update")
		return
	}
	utils.SuccessResponse(ctx, 200, nil, "Marked as invited")
}

// AdminFeedback godoc — GET /admin/beta/feedback
func (c *BetaController) AdminFeedback(ctx *gin.Context) {
	rows := []map[string]interface{}{}
	if err := database.DB.Raw(`SELECT id, email, category, message, page_url, status, created_at FROM feedback ORDER BY created_at DESC LIMIT 500`).Scan(&rows).Error; err != nil {
		utils.InternalServerError(ctx, err, "Failed to load feedback")
		return
	}
	utils.SuccessResponse(ctx, 200, gin.H{"feedback": rows, "total": len(rows)}, "")
}

// AdminFeedbackStatus godoc — PATCH /admin/beta/feedback/:id  {status}
func (c *BetaController) AdminFeedbackStatus(ctx *gin.Context) {
	var req struct {
		Status string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil || (req.Status != "new" && req.Status != "seen" && req.Status != "resolved") {
		utils.BadRequest(ctx, "status must be new, seen or resolved")
		return
	}
	if err := database.DB.Exec(`UPDATE feedback SET status = ? WHERE id = ?`, req.Status, ctx.Param("id")).Error; err != nil {
		utils.InternalServerError(ctx, err, "Failed to update")
		return
	}
	utils.SuccessResponse(ctx, 200, nil, "Updated")
}

// DataHealth godoc — GET /health/data
// Point an uptime monitor (cron-job.org / UptimeRobot) at this. It returns 503 when the scrapers have stopped:
//   - prices older than 45 min while the ESX is open (weekdays 09:30-15:00 EAT), or older than 4 days at any time
//   - no news article stored in the last 36 hours
func (c *BetaController) DataHealth(ctx *gin.Context) {
	var latestPrice, latestArticle *time.Time
	database.DB.Raw(`SELECT max(timestamp) FROM market_prices WHERE ticker_symbol IN (SELECT ticker_symbol FROM company_profiles)`).Scan(&latestPrice)
	database.DB.Raw(`SELECT max(created_at) FROM articles`).Scan(&latestArticle)

	now := time.Now()
	loc, err := time.LoadLocation("Africa/Addis_Ababa")
	if err != nil {
		loc = time.FixedZone("EAT", 3*60*60)
	}
	eat := now.In(loc)
	weekday := eat.Weekday() != time.Saturday && eat.Weekday() != time.Sunday
	marketOpen := weekday && utils.IsEATMarketOpen(now)

	problems := []string{}
	out := gin.H{"market_open": marketOpen, "checked_at": now.UTC()}
	if latestPrice == nil {
		problems = append(problems, "no prices stored")
	} else {
		age := now.Sub(*latestPrice)
		out["latest_price_at"] = latestPrice.UTC()
		out["price_age_minutes"] = int(age.Minutes())
		if marketOpen && age > 45*time.Minute {
			problems = append(problems, "prices stale during market hours")
		}
		if age > 96*time.Hour {
			problems = append(problems, "prices older than 4 days")
		}
	}
	if latestArticle == nil {
		problems = append(problems, "no articles stored")
	} else {
		age := now.Sub(*latestArticle)
		out["latest_article_at"] = latestArticle.UTC()
		out["article_age_hours"] = int(age.Hours())
		if age > 36*time.Hour {
			problems = append(problems, "news stale (>36h)")
		}
	}
	out["problems"] = problems
	if len(problems) > 0 {
		out["status"] = "degraded"
		ctx.JSON(http.StatusServiceUnavailable, out)
		return
	}
	out["status"] = "ok"
	ctx.JSON(http.StatusOK, out)
}
