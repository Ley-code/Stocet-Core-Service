package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupChatRoutes(router *gin.RouterGroup, controller *controllers.ChatController) {
	chat := router.Group("/chat", middleware.AuthMiddleware())
	{
		chat.GET("/channels", controller.ListChannels)
		chat.GET("/channels/:id/messages", controller.GetMessages)
		chat.POST("/channels/:id/messages", controller.SendMessage)
		chat.PATCH("/message/:id", controller.EditMessage)
		chat.DELETE("/message/:id", controller.DeleteMessage)
		chat.POST("/messages/:id/report", controller.ReportMessage)
		chat.POST("/channels", middleware.RequireAdmin(), controller.CreateChannel)
	}
}
