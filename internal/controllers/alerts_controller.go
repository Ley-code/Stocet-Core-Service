package controllers

import (
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
	userID, _ := ctx.Get("user_id")

	alerts, err := c.service.GetByUserID(userID.(string))
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch alerts")
		return
	}

	utils.SuccessResponse(ctx, 200, alerts, "")
}

func (c *AlertsController) Create(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var alert models.Alert
	if err := ctx.ShouldBindJSON(&alert); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		utils.BadRequest(ctx, "Invalid user ID format")
		return
	}

	alert.UserID = userUUID
	alert.Status = "active"

	if err := c.service.Create(&alert); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, alert, "Alert created successfully")
}

func (c *AlertsController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, _ := ctx.Get("user_id")

	// Get existing alert to verify ownership
	alert, err := c.service.GetByID(id)
	if err != nil {
		utils.NotFound(ctx, "Alert not found")
		return
	}

	if alert.UserID.String() != userID.(string) {
		utils.Forbidden(ctx, "You can only update your own alerts")
		return
	}

	var updateAlert models.Alert
	if err := ctx.ShouldBindJSON(&updateAlert); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		utils.BadRequest(ctx, "Invalid user ID format")
		return
	}

	updateAlert.ID = id
	updateAlert.UserID = userUUID

	if err := c.service.Update(&updateAlert); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, updateAlert, "Alert updated successfully")
}

func (c *AlertsController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	userID, _ := ctx.Get("user_id")

	// Get existing alert to verify ownership
	alert, err := c.service.GetByID(id)
	if err != nil {
		utils.NotFound(ctx, "Alert not found")
		return
	}

	if alert.UserID.String() != userID.(string) {
		utils.Forbidden(ctx, "You can only delete your own alerts")
		return
	}

	if err := c.service.Delete(id); err != nil {
		utils.InternalServerError(ctx, err, "Failed to delete alert")
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "Alert deleted successfully")
}
