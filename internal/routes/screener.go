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
	// Open to guests, rate limited per IP.
	screener := router.Group("/screener", middleware.OptionalAuthMiddleware(), middleware.RateLimitByIP(60, 30))
	{
		screener.GET("", controller.Search)
	}
}
