package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAnalyticsRoutes(router *gin.RouterGroup, controller *controllers.AnalyticsController) {
	analytics := router.Group("/analytics")
	{
		analytics.GET("/dashboard", controller.GetDashboardMetrics)
		analytics.GET("/market-overview", controller.GetMarketOverview)
		analytics.GET("/user-activity/:user_id", middleware.AuthMiddleware(), middleware.RequireAdmin(), controller.GetUserActivity)
		analytics.GET("/platform-health", middleware.AuthMiddleware(), middleware.RequireAdmin(), controller.GetPlatformHealth)
	}
}
