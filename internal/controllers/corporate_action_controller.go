package controllers

import (
	"errors"
	"strings"
	"time"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CorporateActionController struct {
	repo *repositories.CorporateActionRepository
}

func NewCorporateActionController(repo *repositories.CorporateActionRepository) *CorporateActionController {
	return &CorporateActionController{repo: repo}
}

// ── Public ────────────────────────────────────────────────────────────────────

// List godoc — GET /corporate-actions
// Query: ticker, type, from (YYYY-MM-DD), to, limit, skip
func (c *CorporateActionController) List(ctx *gin.Context) {
	f := repositories.CorporateActionFilters{
		Ticker:     ctx.Query("ticker"),
		ActionType: ctx.Query("type"),
		Limit:      utils.ParseInt(ctx.DefaultQuery("limit", "50")),
		Skip:       utils.ParseInt(ctx.DefaultQuery("skip", "0")),
	}
	if s := strings.TrimSpace(ctx.Query("from")); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			f.From = &t
		}
	}
	if s := strings.TrimSpace(ctx.Query("to")); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			f.To = &t
		}
	}

	actions, total, err := c.repo.List(f)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch corporate actions")
		return
	}
	utils.SuccessResponse(ctx, 200, gin.H{
		"corporate_actions": actions,
		"total":             total,
		"limit":             f.Limit,
		"skip":              f.Skip,
	}, "")
}

// ── Admin ─────────────────────────────────────────────────────────────────────

type corporateActionBody struct {
	CompanyProfileID string  `json:"company_profile_id" binding:"required"`
	ActionType       string  `json:"action_type"         binding:"required"`
	Title            string  `json:"title"               binding:"required"`
	EventDate        string  `json:"event_date"          binding:"required"` // YYYY-MM-DD
	ExDate           *string `json:"ex_date"`
	RecordDate       *string `json:"record_date"`
	Description      *string `json:"description"`
	Amount           *float64 `json:"amount"`
	Currency         string  `json:"currency"`
}

func parseDate(s *string) (*time.Time, error) {
	if s == nil || strings.TrimSpace(*s) == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", strings.TrimSpace(*s))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// AdminCreate godoc — POST /admin/corporate-actions
func (c *CorporateActionController) AdminCreate(ctx *gin.Context) {
	var body corporateActionBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	cpID, err := uuid.Parse(body.CompanyProfileID)
	if err != nil {
		utils.BadRequest(ctx, "invalid company_profile_id")
		return
	}

	eventDate, err := time.Parse("2006-01-02", strings.TrimSpace(body.EventDate))
	if err != nil {
		utils.BadRequest(ctx, "event_date must be YYYY-MM-DD")
		return
	}

	exDate, err := parseDate(body.ExDate)
	if err != nil {
		utils.BadRequest(ctx, "ex_date must be YYYY-MM-DD")
		return
	}
	recordDate, err := parseDate(body.RecordDate)
	if err != nil {
		utils.BadRequest(ctx, "record_date must be YYYY-MM-DD")
		return
	}

	currency := strings.TrimSpace(body.Currency)
	if currency == "" {
		currency = "ETB"
	}

	action := &models.CorporateAction{
		CompanyProfileID: cpID.String(),
		ActionType:       models.CorporateActionType(body.ActionType),
		Title:            strings.TrimSpace(body.Title),
		EventDate:        eventDate,
		ExDate:           exDate,
		RecordDate:       recordDate,
		Description:      body.Description,
		Amount:           body.Amount,
		Currency:         currency,
	}

	if err := c.repo.Create(action); err != nil {
		utils.InternalServerError(ctx, err, "Failed to create corporate action")
		return
	}
	utils.SuccessResponse(ctx, 201, action, "Corporate action created")
}

// AdminUpdate godoc — PUT /admin/corporate-actions/:id
func (c *CorporateActionController) AdminUpdate(ctx *gin.Context) {
	existing, err := c.repo.GetByID(ctx.Param("id"))
	if err != nil {
		if errors.Is(err, repositories.ErrCorporateActionNotFound) {
			utils.NotFound(ctx, "Corporate action not found")
			return
		}
		utils.InternalServerError(ctx, err, "Failed to fetch corporate action")
		return
	}

	var body corporateActionBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	cpID, err := uuid.Parse(body.CompanyProfileID)
	if err != nil {
		utils.BadRequest(ctx, "invalid company_profile_id")
		return
	}
	eventDate, err := time.Parse("2006-01-02", strings.TrimSpace(body.EventDate))
	if err != nil {
		utils.BadRequest(ctx, "event_date must be YYYY-MM-DD")
		return
	}
	exDate, err := parseDate(body.ExDate)
	if err != nil {
		utils.BadRequest(ctx, "ex_date must be YYYY-MM-DD")
		return
	}
	recordDate, err := parseDate(body.RecordDate)
	if err != nil {
		utils.BadRequest(ctx, "record_date must be YYYY-MM-DD")
		return
	}

	currency := strings.TrimSpace(body.Currency)
	if currency == "" {
		currency = "ETB"
	}

	existing.CompanyProfileID = cpID.String()
	existing.ActionType = models.CorporateActionType(body.ActionType)
	existing.Title = strings.TrimSpace(body.Title)
	existing.EventDate = eventDate
	existing.ExDate = exDate
	existing.RecordDate = recordDate
	existing.Description = body.Description
	existing.Amount = body.Amount
	existing.Currency = currency

	if err := c.repo.Update(existing); err != nil {
		utils.InternalServerError(ctx, err, "Failed to update corporate action")
		return
	}
	utils.SuccessResponse(ctx, 200, existing, "Corporate action updated")
}

// AdminDelete godoc — DELETE /admin/corporate-actions/:id
func (c *CorporateActionController) AdminDelete(ctx *gin.Context) {
	if err := c.repo.Delete(ctx.Param("id")); err != nil {
		if errors.Is(err, repositories.ErrCorporateActionNotFound) {
			utils.NotFound(ctx, "Corporate action not found")
			return
		}
		utils.InternalServerError(ctx, err, "Failed to delete corporate action")
		return
	}
	utils.SuccessResponse(ctx, 200, nil, "Corporate action deleted")
}
