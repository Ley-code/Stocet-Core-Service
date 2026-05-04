package controllers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// AIController exposes the /api/v1/ai/chat endpoint.
//
// The controller does no LLM-specific logic itself: it validates the request,
// delegates to AIService, and maps service errors to graceful HTTP responses.
type AIController struct {
	svc *services.AIService
}

// NewAIController is the canonical constructor used in main.go.
func NewAIController(svc *services.AIService) *AIController {
	return &AIController{svc: svc}
}

// chatRequest is the wire payload for POST /ai/chat.
type chatRequest struct {
	Messages []services.ChatMessage `json:"messages"`
}

// gracefulErrorMessage is the canonical user-facing error per the plan.
const gracefulErrorMessage = "The data engine is currently experiencing high load. Please try again in a moment."

// per-message hard cap so we don't ship absurd payloads to Gemini.
const maxMessageContentChars = 8000

// Chat handles POST /ai/chat. It assumes AuthMiddleware has already gated
// the route so we don't recheck authentication here.
func (c *AIController) Chat(ctx *gin.Context) {
	if c.svc == nil || !c.svc.Configured() {
		respondGraceful(ctx)
		return
	}

	var req chatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "Invalid request body")
		return
	}
	if len(req.Messages) == 0 {
		utils.BadRequest(ctx, "messages must contain at least one entry")
		return
	}

	for i, m := range req.Messages {
		role := strings.ToLower(strings.TrimSpace(m.Role))
		if role != "user" && role != "assistant" {
			utils.BadRequest(ctx, "messages[].role must be 'user' or 'assistant'")
			return
		}
		req.Messages[i].Role = role
		req.Messages[i].Content = strings.TrimSpace(m.Content)
		if req.Messages[i].Content == "" {
			utils.BadRequest(ctx, "messages[].content must not be empty")
			return
		}
		if len(req.Messages[i].Content) > maxMessageContentChars {
			utils.BadRequest(ctx, "messages[].content exceeds maximum length")
			return
		}
	}
	if req.Messages[len(req.Messages)-1].Role != "user" {
		utils.BadRequest(ctx, "the last message must be from the user")
		return
	}

	// Bound the upstream call so we can fail gracefully if Gemini stalls.
	callCtx, cancel := context.WithTimeout(ctx.Request.Context(), 60*time.Second)
	defer cancel()

	answer, err := c.svc.Chat(callCtx, req.Messages)
	if err != nil {
		if errors.Is(err, services.ErrAIDataEngineUnavailable) {
			respondGraceful(ctx)
			return
		}
		respondGraceful(ctx)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"message": gin.H{
			"role":    "assistant",
			"content": answer,
		},
	}, "")
}

// respondGraceful renders the canonical 503 body the frontend renders inline.
func respondGraceful(ctx *gin.Context) {
	ctx.JSON(http.StatusServiceUnavailable, utils.APIResponse{
		Success: false,
		Error:   gracefulErrorMessage,
		Message: gracefulErrorMessage,
	})
}
