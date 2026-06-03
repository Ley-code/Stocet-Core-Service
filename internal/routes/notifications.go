package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupNotificationsRoutes(
	router *gin.RouterGroup,
	controller *controllers.NotificationsController,
	prefController *controllers.NotificationPreferenceController,
) {
	// Internal endpoint — called by the News Aggregation Service.
	router.POST("/internal/trigger-alerts", middleware.InternalSecretMiddleware(), controller.HandleTriggerAlerts)

	// User-facing notification endpoints (require authentication).
	notifications := router.Group("/notifications", middleware.AuthMiddleware())
	{
		notifications.GET("", controller.GetAllNotifications)
		notifications.GET("/unread", controller.GetUnreadNotifications)
		notifications.PATCH("/:id/read", controller.MarkNotificationRead)

		// Notification preferences
		notifications.GET("/preferences", prefController.Get)
		notifications.PUT("/preferences", prefController.Update)
	}
}
