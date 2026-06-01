package middleware

import (
	"stocet-core-service/internal/config"
	"stocet-core-service/internal/services"

	"github.com/gin-gonic/gin"
)

// RequireSubscription aborts with 402 if the authenticated user's effective subscription
// level (resolved from DB, including org plan) is below minLevel. Admins always pass.
func RequireSubscription(subSvc *services.SubscriptionService, minLevel string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, _ := c.Get("user_role"); role == config.ADMIN_ROLE {
			c.Next()
			return
		}

		userID, _ := c.Get("user_id")
		userIDStr, _ := userID.(string)

		ok, effective, err := subSvc.CanAccess(userIDStr, minLevel)
		if err != nil {
			c.JSON(500, gin.H{"success": false, "error": "failed to verify subscription"})
			c.Abort()
			return
		}
		if !ok {
			c.JSON(402, gin.H{
				"success":  false,
				"error":    "subscription_required",
				"message":  "This feature requires a " + minLevel + " subscription. Upgrade to unlock it.",
				"required": minLevel,
				"current":  effective,
			})
			c.Abort()
			return
		}
		c.Set("effective_subscription", effective)
		c.Next()
	}
}
