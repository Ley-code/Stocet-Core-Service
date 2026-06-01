package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupScreenerRoutes(router *gin.RouterGroup, controller *controllers.ScreenerController) {
	if controller == nil {
		return
	}
	// Auth required — data is not sensitive but keeps anonymous traffic out.
	screener := router.Group("/screener", middleware.AuthMiddleware())
	{
		screener.GET("", controller.Search)
	}
}
