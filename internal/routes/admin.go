package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"
	"stocet-core-service/internal/repositories"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(
	router *gin.RouterGroup,
	adminController *controllers.AdminController,
	inviteCodesController *controllers.InviteCodesAdminController,
	orgController *controllers.OrganizationController,
	caController *controllers.CorporateActionController,
	auditController *controllers.AuditLogController,
	analyticsController *controllers.AdminAnalyticsController,
	fundamentalsController *controllers.FundamentalsController,
	embeddingsAdmin *controllers.EmbeddingsAdminController,
	auditRepo *repositories.AuditLogRepository,
) {
	admin := router.Group("/admin",
		middleware.AuthMiddleware(),
		middleware.RequireAdmin(),
		middleware.AuditMiddleware(auditRepo),
	)
	{
		// Users
		admin.GET("/users", adminController.ListUsers)
		admin.PUT("/users/:id", adminController.UpdateUser)
		admin.DELETE("/users/:id", adminController.DeleteUser)

		// Invite codes
		admin.GET("/invite-codes", inviteCodesController.ListInviteCodes)
		admin.POST("/invite-codes", inviteCodesController.CreateInviteCode)
		admin.PATCH("/invite-codes/:id/revoke", inviteCodesController.RevokeInviteCode)
		admin.DELETE("/invite-codes/:id", inviteCodesController.DeleteInviteCode)

		// Audit log (read-only — GET is skipped by AuditMiddleware)
		admin.GET("/audit-log", auditController.List)

		// Usage analytics
		admin.GET("/analytics", analyticsController.Get)
	}
	RegisterOrganizationRoutes(admin, orgController)
	RegisterCorporateActionAdminRoutes(admin, caController)
	if fundamentalsController != nil {
		RegisterFundamentalsRoutes(admin, fundamentalsController)
	}
	if embeddingsAdmin != nil {
		admin.POST("/ai/embeddings/backfill", embeddingsAdmin.RunBackfill)
		admin.POST("/ai/sentiment/backfill", embeddingsAdmin.RunSentimentBackfill)
	}
}
