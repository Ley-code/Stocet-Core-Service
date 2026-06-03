package controllers

import (
	"time"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EarningsController struct {
	repo *repositories.EarningsRepository
}

func NewEarningsController(repo *repositories.EarningsRepository) *EarningsController {
	return &EarningsController{repo: repo}
}

func (c *EarningsController) List(ctx *gin.Context) {
	var from, to *time.Time
	if s := ctx.Query("from"); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			from = &t
		}
	}
	if s := ctx.Query("to"); s != "" {
		if t, err := time.Parse("2006-01-02", s); err == nil {
			to = &t
		}
	}
	ticker := ctx.Query("ticker")
	status := ctx.Query("status")

	items, err := c.repo.List(from, to, ticker, status)
	if err != nil {
		utils.InternalServerError(ctx, err, "failed to list earnings")
		return
	}
	if items == nil {
		items = []models.EarningsDate{}
	}
	utils.SuccessResponse(ctx, 200, gin.H{"earnings": items}, "")
}

func (c *EarningsController) Create(ctx *gin.Context) {
	var body models.EarningsDate
	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(ctx, "invalid request body")
		return
	}
	if err := c.repo.Create(&body); err != nil {
		utils.InternalServerError(ctx, err, "failed to create earnings date")
		return
	}
	utils.SuccessResponse(ctx, 201, body, "")
}

func (c *EarningsController) Update(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.BadRequest(ctx, "invalid id")
		return
	}
	var body map[string]interface{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(ctx, "invalid request body")
		return
	}
	allowed := map[string]bool{
		"report_period":    true,
		"report_type":      true,
		"expected_date":    true,
		"actual_date":      true,
		"status":           true,
		"notes":            true,
		"eps_estimate":     true,
		"eps_actual":       true,
		"revenue_estimate": true,
		"revenue_actual":   true,
	}
	updates := make(map[string]interface{})
	for k, v := range body {
		if allowed[k] {
			updates[k] = v
		}
	}
	if err := c.repo.Update(id, updates); err != nil {
		utils.InternalServerError(ctx, err, "failed to update earnings date")
		return
	}
	utils.SuccessResponse(ctx, 200, gin.H{"id": id}, "")
}

func (c *EarningsController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.BadRequest(ctx, "invalid id")
		return
	}
	if err := c.repo.Delete(id); err != nil {
		utils.InternalServerError(ctx, err, "failed to delete earnings date")
		return
	}
	utils.SuccessResponse(ctx, 200, gin.H{"deleted": true}, "")
}
