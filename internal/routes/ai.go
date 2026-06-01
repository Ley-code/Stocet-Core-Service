package routes

import (
	"stocet-core-service/internal/config"
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"
	"stocet-core-service/internal/services"

	"github.com/gin-gonic/gin"
)

// SetupAIRoutes mounts the AI Data Assistant endpoints. All routes require
// authentication and a premium effective subscription tier.
func SetupAIRoutes(router *gin.RouterGroup, controller *controllers.AIController, subSvc *services.SubscriptionService) {
	if controller == nil {
		return
	}
	ai := router.Group("/ai", middleware.AuthMiddleware())
	if subSvc != nil {
		ai.Use(middleware.RequireSubscription(subSvc, config.SUBSCRIPTION_PREMIUM))
	}
	{
		ai.GET("/usage", controller.Usage)
		ai.POST("/chat", controller.Chat)
		ai.POST("/chat/stream", controller.ChatStream)
	}
}
