package controllers

import (
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	service *services.DashboardService
}

func NewDashboardController(service *services.DashboardService) *DashboardController {
	return &DashboardController{service: service}
}

func (c *DashboardController) GetDashboard(ctx *gin.Context) {
	// Guests get market data without personalised alerts (empty userID).
	userID, _ := utils.GetUserID(ctx)

	data, err := c.service.GetDashboardData(userID)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch dashboard data")
		return
	}

	utils.SuccessResponse(ctx, 200, data, "")
}
