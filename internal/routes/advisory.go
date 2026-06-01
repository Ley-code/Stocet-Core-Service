package routes

import (
	"stocet-core-service/internal/config"
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"
	"stocet-core-service/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupAdvisoryRoutes(router *gin.RouterGroup, controller *controllers.AdvisoryController, subSvc *services.SubscriptionService) {
	if controller == nil {
		return
	}
	advisory := router.Group("/advisory", middleware.AuthMiddleware())
	if subSvc != nil {
		advisory.Use(middleware.RequireSubscription(subSvc, config.SUBSCRIPTION_PREMIUM))
	}
	{
		advisory.GET("/reports", controller.List)
		advisory.POST("/reports", controller.Create)
		advisory.DELETE("/reports/:id", controller.Delete)
	}
}
