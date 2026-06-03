package controllers

import (
	"log"

	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type NotificationsController struct {
	service *services.NotificationsService
}

func NewNotificationsController(service *services.NotificationsService) *NotificationsController {
	return &NotificationsController{service: service}
}

// HandleTriggerAlerts is an internal endpoint called by the News Aggregation
// Service after new market prices are ingested. It matches active alerts
// against the latest prices and creates notifications for triggered alerts.
func (ctrl *NotificationsController) HandleTriggerAlerts(c *gin.Context) {
	count, err := ctrl.service.TriggerAlerts()
	if err != nil {
		utils.InternalServerError(c, err, "Failed to process alerts")
		return
	}

	log.Printf("Trigger alerts complete: %d alerts matched", count)
	utils.SuccessResponse(c, 200, gin.H{"triggered": count}, "Alerts processed successfully")
}

// GetUnreadNotifications returns all unread notifications for the authenticated user.
func (ctrl *NotificationsController) GetUnreadNotifications(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	notifications, err := ctrl.service.GetUnread(userID)
	if err != nil {
		utils.InternalServerError(c, err, "Failed to fetch unread notifications")
		return
	}

	utils.SuccessResponse(c, 200, notifications, "")
}

// MarkNotificationRead marks a single notification as read.
func (ctrl *NotificationsController) MarkNotificationRead(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		utils.Unauthorized(c, "User not authenticated")
		return
	}
	id := c.Param("id")

	if err := ctrl.service.MarkRead(id, userID); err != nil {
		utils.NotFound(c, "Notification not found")
		return
	}

	utils.SuccessResponse(c, 200, nil, "Notification marked as read")
}

// GetAllNotifications returns all notifications for the authenticated user,
// ordered by created_at DESC.
func (ctrl *NotificationsController) GetAllNotifications(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		utils.Unauthorized(c, "User not authenticated")
		return
	}

	notifications, err := ctrl.service.GetAll(userID)
	if err != nil {
		utils.InternalServerError(c, err, "Failed to fetch notifications")
		return
	}

	utils.SuccessResponse(c, 200, notifications, "")
}
