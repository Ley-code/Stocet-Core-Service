package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupDashboardRoutes(router *gin.RouterGroup, controller *controllers.DashboardController) {
	dashboard := router.Group("/dashboard", middleware.AuthMiddleware())
	{
		dashboard.GET("", controller.GetDashboard)
	}
}
