package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupCapitalServiceProviderRoutes(router *gin.RouterGroup, controller *controllers.CapitalServiceProvidersController) {
	g := router.Group("/capital-service-providers", middleware.AuthMiddleware())
	{
		g.GET("", controller.List)
		g.GET("/:id", controller.GetByID)
		g.POST("", middleware.RequireAdmin(), controller.Create)
		g.PUT("/:id", middleware.RequireAdmin(), controller.Update)
		g.DELETE("/:id", middleware.RequireAdmin(), controller.Delete)
	}
}
