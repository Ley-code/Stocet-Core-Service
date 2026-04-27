package controllers

import (
	"stocet-core-service/internal/models"
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
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
	userID, _ := ctx.Get("user_id")
	channelID := ctx.Param("id")

	var req struct {
		Content  string  `json:"content" binding:"required"`
		ParentID *string `json:"parent_id,omitempty"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	message, err := c.service.CreateMessage(channelID, userID.(string), req.Content, req.ParentID)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, message, "Message sent successfully")
}

func (c *ChatController) EditMessage(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	messageID := ctx.Param("id")

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	message, err := c.service.EditMessage(messageID, userID.(string), req.Content)
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
	userID, _ := ctx.Get("user_id")
	messageID := ctx.Param("id")

	err := c.service.DeleteMessage(messageID, userID.(string))
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
