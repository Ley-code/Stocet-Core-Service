package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupWatchlistRoutes(router *gin.RouterGroup, controller *controllers.WatchlistController) {
	watchlist := router.Group("/watchlist", middleware.AuthMiddleware())
	{
		watchlist.GET("", controller.List)
		watchlist.POST("", controller.Add)
		watchlist.DELETE("/:id", controller.Remove)
	}
}
