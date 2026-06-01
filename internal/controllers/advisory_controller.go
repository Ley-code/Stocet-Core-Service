package controllers

import (
	"encoding/json"
	"errors"
	"strings"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AdvisoryController struct {
	repo *repositories.AdvisoryReportRepository
}

func NewAdvisoryController(repo *repositories.AdvisoryReportRepository) *AdvisoryController {
	return &AdvisoryController{repo: repo}
}

// List godoc — GET /advisory/reports
func (c *AdvisoryController) List(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	reports, err := c.repo.ListByUser(userID.(string))
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch advisory reports")
		return
	}
	utils.SuccessResponse(ctx, 200, gin.H{"reports": reports, "total": len(reports)}, "")
}

// Create godoc — POST /advisory/reports
// Body: { title: string, profile: object, result: string }
func (c *AdvisoryController) Create(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var body struct {
		Title   string                 `json:"title"`
		Profile map[string]interface{} `json:"profile"`
		Result  string                 `json:"result"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}
	if strings.TrimSpace(body.Result) == "" {
		utils.BadRequest(ctx, "result is required")
		return
	}

	profileBytes, _ := json.Marshal(body.Profile)

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		utils.BadRequest(ctx, "invalid user id")
		return
	}

	report := &models.AdvisoryReport{
		UserID:  uid,
		Title:   strings.TrimSpace(body.Title),
		Profile: datatypes.JSON(profileBytes),
		Result:  body.Result,
	}

	if err := c.repo.Create(report); err != nil {
		utils.InternalServerError(ctx, err, "Failed to save advisory report")
		return
	}
	utils.SuccessResponse(ctx, 201, report, "Report saved")
}

// Delete godoc — DELETE /advisory/reports/:id
func (c *AdvisoryController) Delete(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	if err := c.repo.Delete(userID.(string), ctx.Param("id")); err != nil {
		if errors.Is(err, repositories.ErrAdvisoryReportNotFound) {
			utils.NotFound(ctx, "Report not found")
			return
		}
		utils.InternalServerError(ctx, err, "Failed to delete report")
		return
	}
	utils.SuccessResponse(ctx, 200, nil, "Report deleted")
}
