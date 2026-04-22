package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupNotificationsRoutes(router *gin.RouterGroup, controller *controllers.NotificationsController) {
	// Internal endpoint — called by the News Aggregation Service.
	// No auth middleware; secure via network/API-key in production.
	router.POST("/internal/trigger-alerts", controller.HandleTriggerAlerts)

	// User-facing notification endpoints (require authentication).
	notifications := router.Group("/notifications", middleware.AuthMiddleware())
	{
		notifications.GET("", controller.GetAllNotifications)
		notifications.GET("/unread", controller.GetUnreadNotifications)
		notifications.PATCH("/:id/read", controller.MarkNotificationRead)
	}
}
