package controllers

import (
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuditLogController struct {
	repo *repositories.AuditLogRepository
}

func NewAuditLogController(repo *repositories.AuditLogRepository) *AuditLogController {
	return &AuditLogController{repo: repo}
}

// List godoc — GET /admin/audit-log
// Query: actor_id, method, limit, skip
func (c *AuditLogController) List(ctx *gin.Context) {
	f := repositories.AuditLogFilters{
		ActorID: ctx.Query("actor_id"),
		Method:  ctx.Query("method"),
		Limit:   utils.ParseInt(ctx.DefaultQuery("limit", "50")),
		Skip:    utils.ParseInt(ctx.DefaultQuery("skip", "0")),
	}

	logs, total, err := c.repo.List(f)
	if err != nil {
		utils.InternalServerError(ctx, err, "Failed to fetch audit log")
		return
	}
	utils.SuccessResponse(ctx, 200, gin.H{
		"logs":  logs,
		"total": total,
		"limit": f.Limit,
		"skip":  f.Skip,
	}, "")
}
