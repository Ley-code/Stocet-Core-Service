package controllers

import (
	"log"
	"net/http"

	"stocet-core-service/internal/config"
	"stocet-core-service/internal/services/ai"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// EmbeddingsAdminController exposes admin-only hooks for vector embedding
// maintenance (e.g. one-shot backfill on Lambda where the periodic worker
// does not run).
type EmbeddingsAdminController struct {
	embedder *ai.Embedder
}

// NewEmbeddingsAdminController constructs a controller. embedder may be nil
// (handler returns 503).
func NewEmbeddingsAdminController(embedder *ai.Embedder) *EmbeddingsAdminController {
	return &EmbeddingsAdminController{embedder: embedder}
}

// RunBackfill handles POST /admin/ai/embeddings/backfill?batch=25
// Triggers a single RunBackfillOnce pass (one batch from corporate_news and
// one from company_profiles). Requires admin auth (route registration).
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
