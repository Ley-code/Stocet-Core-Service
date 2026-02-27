package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupNewsRoutes(router *gin.RouterGroup, controller *controllers.NewsController) {
	news := router.Group("/news", middleware.AuthMiddleware())
	{
		news.GET("", controller.GetNews)
		news.DELETE("/:id", middleware.RequireAdmin(), controller.DeleteNews)
	}
}
