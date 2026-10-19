package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupCapitalServiceProviderRoutes(router *gin.RouterGroup, controller *controllers.CapitalServiceProvidersController) {
	open := router.Group("/capital-service-providers", middleware.OptionalAuthMiddleware(), middleware.RateLimitByIP(60, 30))
	{
		open.GET("", controller.List)
		open.GET("/:id", controller.GetByID)
	}

	g := router.Group("/capital-service-providers", middleware.AuthMiddleware())
	{
		g.POST("", middleware.RequireAdmin(), controller.Create)
		g.PUT("/:id", middleware.RequireAdmin(), controller.Update)
		g.DELETE("/:id", middleware.RequireAdmin(), controller.Delete)
	}
}
