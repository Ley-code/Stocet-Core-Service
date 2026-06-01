package controllers

import (
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationPreferenceController struct {
	repo *repositories.NotificationPreferenceRepository
}

func NewNotificationPreferenceController(
	repo *repositories.NotificationPreferenceRepository,
) *NotificationPreferenceController {
	return &NotificationPreferenceController{repo: repo}
}

// Get godoc — GET /notifications/preferences
func (c *NotificationPreferenceController) Get(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	pref, err := c.repo.Get(userID.(string))
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch notification preferences")
		return
	}
	utils.SuccessResponse(ctx, 200, pref, "")
}

// Update godoc — PUT /notifications/preferences
func (c *NotificationPreferenceController) Update(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	uid, err := uuid.Parse(userID.(string))
	if err != nil {
		utils.BadRequest(ctx, "invalid user id")
		return
	}

	// Load existing first (creates defaults if missing)
	pref, err := c.repo.Get(userID.(string))
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch preferences")
		return
	}

	var body struct {
		PriceAlertPush  *bool `json:"price_alert_push"`
		PriceAlertEmail *bool `json:"price_alert_email"`
		NewsAlertPush   *bool `json:"news_alert_push"`
		NewsAlertEmail  *bool `json:"news_alert_email"`
		InAppEnabled    *bool `json:"in_app_enabled"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	pref.UserID = uid
	if body.PriceAlertPush  != nil { pref.PriceAlertPush  = *body.PriceAlertPush }
	if body.PriceAlertEmail != nil { pref.PriceAlertEmail = *body.PriceAlertEmail }
	if body.NewsAlertPush   != nil { pref.NewsAlertPush   = *body.NewsAlertPush }
	if body.NewsAlertEmail  != nil { pref.NewsAlertEmail  = *body.NewsAlertEmail }
	if body.InAppEnabled    != nil { pref.InAppEnabled    = *body.InAppEnabled }

	if err := c.repo.Save(pref); err != nil {
		utils.InternalServerError(ctx, err, "Failed to save preferences")
		return
	}
	utils.SuccessResponse(ctx, 200, pref, "Preferences updated")
}
