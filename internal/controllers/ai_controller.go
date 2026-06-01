package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"stocet-core-service/internal/config"
	"stocet-core-service/internal/services"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// AIController exposes the /api/v1/ai/chat endpoint.
//
// The controller does no LLM-specific logic itself: it validates the request,
// delegates to AIService, and maps service errors to graceful HTTP responses.
type AIController struct {
	svc      *services.AIService
	quotaSvc *services.AIQuotaService
}

// NewAIController is the canonical constructor used in main.go.
func NewAIController(svc *services.AIService, quotaSvc *services.AIQuotaService) *AIController {
	return &AIController{svc: svc, quotaSvc: quotaSvc}
}

func (c *AIController) isAdmin(ctx *gin.Context) bool {
	role, _ := ctx.Get("user_role")
	return role == config.ADMIN_ROLE
}

func (c *AIController) enforceQuota(ctx *gin.Context) bool {
	if c.quotaSvc == nil {
		return true
	}
	userID, _ := ctx.Get("user_id")
	if err := c.quotaSvc.CheckAndIncrement(userID.(string), c.isAdmin(ctx)); err != nil {
		if errors.Is(err, services.ErrAIQuotaExceeded) {
			used, limit, _ := c.quotaSvc.GetUsage(userID.(string), false)
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "ai_quota_exceeded",
				"message": "Daily AI request limit reached. Try again tomorrow or upgrade your plan.",
				"used":    used,
				"limit":   limit,
			})
			return false
		}
		utils.InternalServerError(ctx, err, "Failed to check AI quota")
		return false
	}
	return true
}

// Usage handles GET /ai/usage — today's request count and daily limit.
func (c *AIController) Usage(ctx *gin.Context) {
	if c.quotaSvc == nil {
		utils.SuccessResponse(ctx, http.StatusOK, gin.H{"used": 0, "limit": 0}, "")
		return
	}
	userID, _ := ctx.Get("user_id")
	used, limit, err := c.quotaSvc.GetUsage(userID.(string), c.isAdmin(ctx))
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch AI usage")
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"used":    used,
		"limit":   limit,
		"unlimited": limit < 0,
	}, "")
}

// chatRequest is the wire payload for POST /ai/chat.
// Mode "advisory" switches to the portfolio guidance system prompt.
type chatRequest struct {
	Messages []services.ChatMessage `json:"messages"`
	Mode     string                 `json:"mode,omitempty"` // "" | "advisory"
}

// gracefulErrorMessage is the canonical user-facing error per the plan.
const gracefulErrorMessage = "The data engine is currently experiencing high load. Please try again in a moment."

// per-message hard cap so we don't ship absurd payloads to Gemini.
const maxMessageContentChars = 8000

// Chat handles POST /ai/chat. It assumes AuthMiddleware has already gated
// the route so we don't recheck authentication here.
func (c *AIController) Chat(ctx *gin.Context) {
	if c.svc == nil {
		log.Printf("[ai][chat] 503 client_ip=%s reason=ai_service_nil", ctx.ClientIP())
		respondGraceful(ctx)
		return
	}
	if !c.svc.Configured() {
		log.Printf("[ai][chat] 503 client_ip=%s reason=gemini_not_configured (set GEMINI_API_KEY in Lambda/env)", ctx.ClientIP())
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
	if !c.enforceQuota(ctx) {
		return
	}

	msgCount := len(req.Messages)
	totalChars := 0
	for _, m := range req.Messages {
		totalChars += len(m.Content)
	}
	log.Printf("[ai][chat] start client_ip=%s messages=%d total_content_chars=%d", ctx.ClientIP(), msgCount, totalChars)

	// Bound the upstream call so we can fail gracefully if Gemini stalls.
	callCtx, cancel := context.WithTimeout(ctx.Request.Context(), 60*time.Second)
	defer cancel()

	answer, err := c.svc.Chat(callCtx, req.Messages, req.Mode)
	if err != nil {
		engine := errors.Is(err, services.ErrAIDataEngineUnavailable)
		log.Printf("[ai][chat] 503 client_ip=%s messages=%d err=%q is_engine_unavailable=%v", ctx.ClientIP(), msgCount, err.Error(), engine)
		respondGraceful(ctx)
		return
	}

	log.Printf("[ai][chat] ok client_ip=%s messages=%d answer_chars=%d", ctx.ClientIP(), msgCount, len(answer))
	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"message": gin.H{
			"role":    "assistant",
			"content": answer,
		},
	}, "")
}

// ChatStream handles POST /ai/chat/stream. It validates the request identically
// to Chat, then emits SSE events: tool progress + the final text response.
func (c *AIController) ChatStream(ctx *gin.Context) {
	if c.svc == nil || !c.svc.Configured() {
		log.Printf("[ai][stream] 503 client_ip=%s reason=not_configured", ctx.ClientIP())
		ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": gracefulErrorMessage})
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
	if !c.enforceQuota(ctx) {
		return
	}

	flusher, ok := ctx.Writer.(http.Flusher)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "streaming not supported"})
		return
	}

	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("X-Accel-Buffering", "no")
	ctx.Header("Connection", "keep-alive")

	log.Printf("[ai][stream] start client_ip=%s messages=%d", ctx.ClientIP(), len(req.Messages))

	callCtx, cancel := context.WithTimeout(ctx.Request.Context(), 90*time.Second)
	defer cancel()

	w := ctx.Writer
	emit := func(event services.StreamEvent) {
		data, _ := json.Marshal(event)
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}

	if err := c.svc.ChatStream(callCtx, req.Messages, req.Mode, emit); err != nil {
		log.Printf("[ai][stream] error client_ip=%s err=%v", ctx.ClientIP(), err)
	}

	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}

// respondGraceful renders the canonical 503 body the frontend renders inline.
func respondGraceful(ctx *gin.Context) {
	ctx.JSON(http.StatusServiceUnavailable, utils.APIResponse{
		Success: false,
		Error:   gracefulErrorMessage,
		Message: gracefulErrorMessage,
	})
}
