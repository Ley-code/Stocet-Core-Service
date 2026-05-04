package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAIRoutes mounts the AI Data Assistant endpoints. All routes require
// the standard AuthMiddleware (logged-in users only).
func SetupAIRoutes(router *gin.RouterGroup, controller *controllers.AIController) {
	if controller == nil {
		return
	}
	ai := router.Group("/ai", middleware.AuthMiddleware())
	{
		ai.POST("/chat", controller.Chat)
	}
}
