package controllers

import (
	"context"
	"log"
	"net/http"

	"stocet-core-service/internal/config"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/services/ai"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// EmbeddingsAdminController exposes admin-only hooks for vector embedding
// maintenance and news sentiment backfill.
type EmbeddingsAdminController struct {
	embedder *ai.Embedder
	gemini   *ai.GeminiClient
	newsRepo *repositories.NewsRepository
}

func NewEmbeddingsAdminController(embedder *ai.Embedder) *EmbeddingsAdminController {
	return &EmbeddingsAdminController{embedder: embedder}
}

// WithSentiment attaches the Gemini client + news repo needed for sentiment backfill.
func (c *EmbeddingsAdminController) WithSentiment(gemini *ai.GeminiClient, newsRepo *repositories.NewsRepository) *EmbeddingsAdminController {
	c.gemini = gemini
	c.newsRepo = newsRepo
	return c
}

// RunBackfill handles POST /admin/ai/embeddings/backfill?batch=25
func (c *EmbeddingsAdminController) RunBackfill(ctx *gin.Context) {
	if c == nil || c.embedder == nil || !c.embedder.Configured() {
		utils.ErrorResponse(ctx, http.StatusServiceUnavailable, nil,
			"Embedding service not configured (set GEMINI_API_KEY and GEMINI_EMBED_MODEL)")
		return
	}

	batch := config.AppConfig.AIBackfillBatchSize
	if batch <= 0 {
		batch = 25
	}
	if q := ctx.Query("batch"); q != "" {
		if v := utils.ParseInt(q); v > 0 && v <= 500 {
			batch = v
		}
	}

	stats := ai.RunBackfillOnce(ctx.Request.Context(), c.embedder, batch)
	log.Printf("[ai] admin backfill batch=%d corporate_news=%d company_profiles=%d",
		batch, stats.CorporateNewsRows, stats.CompanyProfilesRows)

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"corporate_news_embedded":   stats.CorporateNewsRows,
		"company_profiles_embedded": stats.CompanyProfilesRows,
		"batch_size":                batch,
	}, "Embedding backfill batch completed")
}

// RunSentimentBackfill handles POST /admin/ai/sentiment/backfill?batch=20
// Classifies one batch of unanalyzed news articles using Gemini.
func (c *EmbeddingsAdminController) RunSentimentBackfill(ctx *gin.Context) {
	if c.gemini == nil || !c.gemini.Configured() {
		utils.ErrorResponse(ctx, http.StatusServiceUnavailable, nil,
			"Sentiment service not configured (set GEMINI_API_KEY)")
		return
	}
	if c.newsRepo == nil {
		utils.ErrorResponse(ctx, http.StatusServiceUnavailable, nil, "News repository not available")
		return
	}

	batch := 20
	if q := ctx.Query("batch"); q != "" {
		if v := utils.ParseInt(q); v > 0 && v <= 100 {
			batch = v
		}
	}

	processed, err := ai.RunSentimentBackfillOnce(context.Background(), c.gemini, c.newsRepo, batch)
	if err != nil {
		utils.InternalServerError(ctx, err, "Sentiment backfill failed")
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, gin.H{
		"processed":  processed,
		"batch_size": batch,
	}, "Sentiment backfill batch completed")
}
