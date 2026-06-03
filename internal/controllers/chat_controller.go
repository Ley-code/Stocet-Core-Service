package controllers

import (
	"time"

	"stocet-core-service/internal/database"
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatController struct {
	service *services.ChatService
}

func NewChatController(service *services.ChatService) *ChatController {
	return &ChatController{service: service}
}

func (c *ChatController) ListChannels(ctx *gin.Context) {
	channels, err := c.service.ListChannels()
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch channels")
		return
	}

	utils.SuccessResponse(ctx, 200, channels, "")
}

func (c *ChatController) GetMessages(ctx *gin.Context) {
	channelID := ctx.Param("id")
	pagination := utils.ParsePagination(ctx, 50, 200)

	messages, total, err := c.service.GetMessages(channelID, pagination.Limit, pagination.Skip)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch messages")
		return
	}

	utils.SuccessResponse(ctx, 200, gin.H{
		"messages": messages,
		"total":    total,
		"limit":    pagination.Limit,
		"skip":     pagination.Skip,
	}, "")
}

func (c *ChatController) SendMessage(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}
	channelID := ctx.Param("id")

	var req struct {
		Content  string  `json:"content" binding:"required"`
		ParentID *string `json:"parent_id,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	message, err := c.service.CreateMessage(channelID, userID, req.Content, req.ParentID)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, message, "Message sent successfully")
}

func (c *ChatController) EditMessage(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}
	messageID := ctx.Param("id")

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	message, err := c.service.EditMessage(messageID, userID, req.Content)
	if err != nil {
		if err.Error() == "forbidden" {
			utils.Forbidden(ctx, "You can only edit your own messages")
			return
		}
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, message, "Message updated successfully")
}

func (c *ChatController) DeleteMessage(ctx *gin.Context) {
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "User not authenticated")
		return
	}
	messageID := ctx.Param("id")

	err := c.service.DeleteMessage(messageID, userID)
	if err != nil {
		if err.Error() == "forbidden" {
			utils.Forbidden(ctx, "You can only delete your own messages")
			return
		}
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, nil, "Message deleted successfully")
}

// Admin endpoints
func (c *ChatController) CreateChannel(ctx *gin.Context) {
	var channel models.Channel
	if err := ctx.ShouldBindJSON(&channel); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	if err := c.service.CreateChannel(&channel); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, channel, "Channel created successfully")
}

// POST /chat/messages/:id/report
func (c *ChatController) ReportMessage(ctx *gin.Context) {
	userIDStr, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "unauthorized")
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.Unauthorized(ctx, "unauthorized")
		return
	}
	messageID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.BadRequest(ctx, "invalid message id")
		return
	}
	var body struct {
		Reason string `json:"reason"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil || body.Reason == "" {
		utils.BadRequest(ctx, "reason is required")
		return
	}
	report := models.ChatReport{
		MessageID:  messageID,
		ReportedBy: userID,
		Reason:     body.Reason,
	}
	if err := database.DB.Create(&report).Error; err != nil {
		utils.InternalServerError(ctx, err, "failed to report")
		return
	}
	utils.SuccessResponse(ctx, 201, report, "")
}

// GET /admin/chat/reports
func (c *ChatController) AdminListReports(ctx *gin.Context) {
	status := ctx.DefaultQuery("status", "pending")
	var reports []models.ChatReport
	q := database.DB.Table("chat_reports cr").
		Select("cr.*, m.content AS message_content, u.full_name AS reporter_name").
		Joins("LEFT JOIN messages m ON m.id = cr.message_id").
		Joins("LEFT JOIN users u ON u.id = cr.reported_by").
		Where("cr.status = ?", status).
		Order("cr.created_at DESC").
		Limit(100)
	if err := q.Scan(&reports).Error; err != nil {
		utils.InternalServerError(ctx, err, "failed to list reports")
		return
	}
	if reports == nil {
		reports = []models.ChatReport{}
	}
	utils.SuccessResponse(ctx, 200, gin.H{"reports": reports}, "")
}

// PATCH /admin/chat/reports/:id
func (c *ChatController) AdminResolveReport(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		utils.BadRequest(ctx, "invalid id")
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.BadRequest(ctx, "invalid body")
		return
	}
	adminIDStr, _ := utils.GetUserID(ctx)
	adminID, _ := uuid.Parse(adminIDStr)
	now := time.Now()
	if err := database.DB.Model(&models.ChatReport{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": body.Status, "reviewed_by": adminID, "reviewed_at": now,
	}).Error; err != nil {
		utils.InternalServerError(ctx, err, "failed to update")
		return
	}
	utils.SuccessResponse(ctx, 200, gin.H{"updated": true}, "")
}

// POST /admin/chat/mute
func (c *ChatController) AdminMuteUser(ctx *gin.Context) {
	adminIDStr, ok := utils.GetUserID(ctx)
	if !ok {
		utils.Unauthorized(ctx, "unauthorized")
		return
	}
	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		utils.Unauthorized(ctx, "unauthorized")
		return
	}
	var body struct {
		UserID    string  `json:"user_id"`
		ChannelID *string `json:"channel_id"`
		Reason    *string `json:"reason"`
		ExpiresAt *string `json:"expires_at"`
	}
	if err := ctx.ShouldBindJSON(&body); err != nil || body.UserID == "" {
		utils.BadRequest(ctx, "user_id required")
		return
	}
	uid, err := uuid.Parse(body.UserID)
	if err != nil {
		utils.BadRequest(ctx, "invalid user_id")
		return
	}
	mute := models.ChatMute{UserID: uid, MutedBy: adminID}
	if body.Reason != nil {
		mute.Reason = body.Reason
	}
	if body.ChannelID != nil {
		cid, err := uuid.Parse(*body.ChannelID)
		if err == nil {
			mute.ChannelID = &cid
		}
	}
	if body.ExpiresAt != nil {
		if t, err := time.Parse(time.RFC3339, *body.ExpiresAt); err == nil {
			mute.ExpiresAt = &t
		}
	}
	if err := database.DB.Create(&mute).Error; err != nil {
		utils.InternalServerError(ctx, err, "failed to mute")
		return
	}
	utils.SuccessResponse(ctx, 201, mute, "")
}

// DELETE /admin/chat/mute/:userId
func (c *ChatController) AdminUnmuteUser(ctx *gin.Context) {
	uid, err := uuid.Parse(ctx.Param("userId"))
	if err != nil {
		utils.BadRequest(ctx, "invalid userId")
		return
	}
	database.DB.Where("user_id = ?", uid).Delete(&models.ChatMute{})
	utils.SuccessResponse(ctx, 200, gin.H{"unmuted": true}, "")
}

// GET /admin/chat/mutes
func (c *ChatController) AdminListMutes(ctx *gin.Context) {
	type muteRow struct {
		models.ChatMute
		UserName string `json:"user_name"`
	}
	var mutes []muteRow
	database.DB.Table("chat_mutes cm").
		Select("cm.*, u.full_name AS user_name").
		Joins("LEFT JOIN users u ON u.id = cm.user_id").
		Where("cm.expires_at IS NULL OR cm.expires_at > NOW()").
		Order("cm.created_at DESC").
		Limit(100).
		Scan(&mutes)
	if mutes == nil {
		mutes = []muteRow{}
	}
	utils.SuccessResponse(ctx, 200, gin.H{"mutes": mutes}, "")
}
