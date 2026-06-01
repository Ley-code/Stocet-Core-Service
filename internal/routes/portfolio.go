package routes

import (
	"stocet-core-service/internal/config"
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"
	"stocet-core-service/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupPortfolioRoutes(router *gin.RouterGroup, controller *controllers.PortfolioController, subSvc *services.SubscriptionService) {
	if controller == nil {
		return
	}
	portfolio := router.Group("/portfolio", middleware.AuthMiddleware())
	if subSvc != nil {
		portfolio.Use(middleware.RequireSubscription(subSvc, config.SUBSCRIPTION_PREMIUM))
	}
	{
		portfolio.GET("/holdings", controller.List)
		portfolio.POST("/holdings", controller.Create)
		portfolio.PUT("/holdings/:id", controller.Update)
		portfolio.DELETE("/holdings/:id", controller.Delete)
		portfolio.GET("/export", controller.Export)
	}
}
