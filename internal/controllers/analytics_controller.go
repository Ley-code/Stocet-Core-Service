package controllers

import (
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type AnalyticsController struct {
	service *services.AnalyticsService
}

func NewAnalyticsController(service *services.AnalyticsService) *AnalyticsController {
	return &AnalyticsController{service: service}
}

func (c *AnalyticsController) GetDashboardMetrics(ctx *gin.Context) {
	metrics, err := c.service.GetDashboardMetrics()
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch dashboard metrics")
		return
	}

	utils.SuccessResponse(ctx, 200, metrics, "")
}

func (c *AnalyticsController) GetMarketOverview(ctx *gin.Context) {
	metrics, err := c.service.GetDashboardMetrics()
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch market overview")
		return
	}

	utils.SuccessResponse(ctx, 200, metrics.MarketSummary, "")
}

func (c *AnalyticsController) GetUserActivity(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	days := 30

	if d := ctx.Query("days"); d != "" {
		if parsed := utils.ParseInt(d); parsed > 0 {
			days = parsed
		}
	}

	activity, err := c.service.GetUserActivity(userID, days)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch user activity")
		return
	}

	utils.SuccessResponse(ctx, 200, activity, "")
}

func (c *AnalyticsController) GetPlatformHealth(ctx *gin.Context) {
	health, err := c.service.GetPlatformHealth()
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch platform health")
		return
	}

	utils.SuccessResponse(ctx, 200, health, "")
}
