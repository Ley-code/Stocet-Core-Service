package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(
	router *gin.RouterGroup,
	adminController *controllers.AdminController,
	inviteCodesController *controllers.InviteCodesAdminController,
	orgController *controllers.OrganizationController,
	fundamentalsController *controllers.FundamentalsController,
	embeddingsAdmin *controllers.EmbeddingsAdminController,
) {
	admin := router.Group("/admin", middleware.AuthMiddleware(), middleware.RequireAdmin())
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
	}
	RegisterOrganizationRoutes(admin, orgController)
	if fundamentalsController != nil {
		RegisterFundamentalsRoutes(admin, fundamentalsController)
	}
	if embeddingsAdmin != nil {
		admin.POST("/ai/embeddings/backfill", embeddingsAdmin.RunBackfill)
	}
}
