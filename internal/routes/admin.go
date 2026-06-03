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
	marketDataAdminController *controllers.MarketDataAdminController,
	chatController *controllers.ChatController,
	auditRepo *repositories.AuditLogRepository,
	earningsController *controllers.EarningsController,
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
		admin.PUT("/users/:id/role", adminController.UpdateUserRole)
		admin.PUT("/users/:id/subscription", adminController.UpdateUserSubscription)
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

		// Market data feed status
		if marketDataAdminController != nil {
			admin.GET("/market-data/status", marketDataAdminController.FeedStatus)
		}

		// Chat moderation
		if chatController != nil {
			adminChat := admin.Group("/chat")
			adminChat.GET("/reports", chatController.AdminListReports)
			adminChat.PATCH("/reports/:id", chatController.AdminResolveReport)
			adminChat.POST("/mute", chatController.AdminMuteUser)
			adminChat.DELETE("/mute/:userId", chatController.AdminUnmuteUser)
			adminChat.GET("/mutes", chatController.AdminListMutes)
		}
	}
	RegisterOrganizationRoutes(admin, orgController)
	RegisterCorporateActionAdminRoutes(admin, caController)
	RegisterEarningsAdminRoutes(admin, earningsController)
	if fundamentalsController != nil {
		RegisterFundamentalsRoutes(admin, fundamentalsController)
	}
	if embeddingsAdmin != nil {
		admin.POST("/ai/embeddings/backfill", embeddingsAdmin.RunBackfill)
		admin.POST("/ai/sentiment/backfill", embeddingsAdmin.RunSentimentBackfill)
	}
}
