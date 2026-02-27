package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupBrokersRoutes(router *gin.RouterGroup, controller *controllers.BrokersController) {
	brokers := router.Group("/brokers", middleware.AuthMiddleware())
	{
		brokers.GET("", controller.List)
		brokers.GET("/:id", controller.GetByID)
		brokers.POST("", middleware.RequireAdmin(), controller.Create)
		brokers.PUT("/:id", middleware.RequireAdmin(), controller.Update)
		brokers.DELETE("/:id", middleware.RequireAdmin(), controller.Delete)
	}
}
