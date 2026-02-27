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
	userID, _ := ctx.Get("user_id")

	data, err := c.service.GetDashboardData(userID.(string))
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch dashboard data")
		return
	}

	utils.SuccessResponse(ctx, 200, data, "")
}
