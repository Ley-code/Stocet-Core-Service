package routes

import (
	"stocet-core-service/internal/controllers"
	"stocet-core-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupBetaRoutes registers the public waitlist + feedback endpoints and the admin review endpoints.
func SetupBetaRoutes(api *gin.RouterGroup, c *controllers.BetaController) {
	api.POST("/waitlist", middleware.RateLimitByIP(5, 5), c.JoinWaitlist)
	api.POST("/feedback", middleware.OptionalAuthMiddleware(), middleware.RateLimitByIP(10, 10), c.SubmitFeedback)

	admin := api.Group("/admin/beta", middleware.AuthMiddleware(), middleware.RequireAdmin())
	{
		admin.GET("/waitlist", c.AdminWaitlist)
		admin.PATCH("/waitlist/:id/invited", c.AdminMarkInvited)
		admin.POST("/waitlist/:id/send-invite", c.AdminSendInvite)
		admin.GET("/codes", c.AdminInviteCodes)
		admin.GET("/feedback", c.AdminFeedback)
		admin.PATCH("/feedback/:id", c.AdminFeedbackStatus)
	}
}
