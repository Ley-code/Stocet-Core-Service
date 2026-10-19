package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupPaperTradeRoutes(router *gin.RouterGroup, controller *controllers.PaperTradeController) {
	trades := router.Group("/market/paper-trades", middleware.AuthMiddleware())
	{
		trades.GET("", controller.List)
		trades.POST("", controller.Create)
		trades.PUT("/:id", controller.Update)
		trades.DELETE("/:id", controller.Delete)
	}
}
