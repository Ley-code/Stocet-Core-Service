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
		Content string `json:"content" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}

	message, err := c.service.CreateMessage(channelID, userID.(string), req.Content)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, message, "Message sent successfully")
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
