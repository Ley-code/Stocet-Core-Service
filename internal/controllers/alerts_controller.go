package controllers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"time"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AlertsController struct {
	service *services.AlertsService
}

func NewAlertsController(service *services.AlertsService) *AlertsController {
	return &AlertsController{service: service}
}

func (c *AlertsController) List(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}
	alerts, err := c.service.GetByUserID(userID, 0)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch alerts")
		return
	}
	utils.SuccessResponse(ctx, 200, alerts, "")
}

func (c *AlertsController) Create(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	var alert models.Alert
	if err := ctx.ShouldBindJSON(&alert); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(ctx, "Invalid user ID format")
		return
	}
	alert.UserID = userUUID
	alert.Status = "active"

	if err := c.service.Create(&alert); err != nil {
		if errors.Is(err, services.ErrAlertLimitReached) {
			ctx.JSON(http.StatusPaymentRequired, gin.H{
				"success": false,
				"error":   "alert_limit_reached",
				"message": err.Error(),
			})
			return
		}
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.SuccessResponse(ctx, 201, alert, "Alert created successfully")
}

func (c *AlertsController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	existing, err := c.service.GetByID(id)
	if err != nil {
		utils.NotFound(ctx, "Alert not found")
		return
	}
	if existing.UserID.String() != userID {
		utils.Forbidden(ctx, "You can only update your own alerts")
		return
	}

	var update models.Alert
	if err := ctx.ShouldBindJSON(&update); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.BadRequest(ctx, "Invalid user ID format")
		return
	}
	update.ID = id
	update.UserID = userUUID

	if err := c.service.Update(&update); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}
	utils.SuccessResponse(ctx, 200, update, "Alert updated successfully")
}

func (c *AlertsController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	existing, err := c.service.GetByID(id)
	if err != nil {
		utils.NotFound(ctx, "Alert not found")
		return
	}
	if existing.UserID.String() != userID {
		utils.Forbidden(ctx, "You can only delete your own alerts")
		return
	}
	if err := c.service.Delete(id); err != nil {
		utils.InternalServerError(ctx, err, "Failed to delete alert")
		return
	}
	utils.SuccessResponse(ctx, 200, nil, "Alert deleted successfully")
}

// Toggle flips an alert between active and disabled.
func (c *AlertsController) Toggle(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	existing, err := c.service.GetByID(id)
	if err != nil {
		utils.NotFound(ctx, "Alert not found")
		return
	}
	if existing.UserID.String() != userID {
		utils.Forbidden(ctx, "You can only update your own alerts")
		return
	}

	role, _ := ctx.Get("user_role")
	roleStr, _ := role.(string)
	newStatus, err := c.service.Toggle(id, userID, services.IsAdminRole(roleStr))
	if err != nil {
		if errors.Is(err, services.ErrAlertLimitReached) {
			ctx.JSON(http.StatusPaymentRequired, gin.H{
				"success": false,
				"error":   "alert_limit_reached",
				"message": err.Error(),
			})
			return
		}
		utils.InternalServerError(ctx, err, "Failed to toggle alert")
		return
	}
	utils.SuccessResponse(ctx, 200, gin.H{"status": newStatus}, "Alert updated")
}

// Export godoc — GET /alerts/export
// Streams a CSV of all alerts for the authenticated user.
func (c *AlertsController) Export(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}
	alerts, err := c.service.GetByUserID(userID, -1)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch alerts")
		return
	}

	filename := fmt.Sprintf("alerts_%s.csv", time.Now().Format("2006-01-02"))
	ctx.Header("Content-Type", "text/csv; charset=utf-8")
	ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	w := csv.NewWriter(ctx.Writer)
	_ = w.Write([]string{
		"Ticker", "Type", "Condition", "Threshold", "Status",
		"Is Recurring", "Cooldown Hours", "Created At", "Last Triggered At",
	})
	for _, a := range alerts {
		lastTriggered := ""
		if a.LastTriggeredAt != nil {
			lastTriggered = a.LastTriggeredAt.Format(time.RFC3339)
		}
		_ = w.Write([]string{
			a.Ticker,
			string(a.AlertType),
			a.Condition,
			fmt.Sprintf("%.4f", a.Threshold),
			string(a.Status),
			fmt.Sprintf("%v", a.IsRecurring),
			fmt.Sprintf("%d", a.CooldownHours),
			a.CreatedAt.Format(time.RFC3339),
			lastTriggered,
		})
	}
	w.Flush()
}

// GetHistory returns the trigger history for one alert.
func (c *AlertsController) GetHistory(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}

	existing, err := c.service.GetByID(id)
	if err != nil {
		utils.NotFound(ctx, "Alert not found")
		return
	}
	if existing.UserID.String() != userID {
		utils.Forbidden(ctx, "Not your alert")
		return
	}

	history, err := c.service.GetHistory(id)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch history")
		return
	}
	utils.SuccessResponse(ctx, 200, history, "")
}
