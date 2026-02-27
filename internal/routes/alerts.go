package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAlertsRoutes(router *gin.RouterGroup, controller *controllers.AlertsController) {
	alerts := router.Group("/alerts", middleware.AuthMiddleware())
	{
		alerts.GET("", controller.List)
		alerts.POST("", controller.Create)
		alerts.PUT("/:id", controller.Update)
		alerts.DELETE("/:id", controller.Delete)
	}
}
