package middleware

import (
	"stocet-core-service/internal/config"

	"github.com/gin-gonic/gin"
)

// subscriptionOrder maps status strings to a numeric level for comparison.
var subscriptionOrder = map[string]int{
	config.SUBSCRIPTION_FREE:    0,
	config.SUBSCRIPTION_PAID:    1,
	config.SUBSCRIPTION_PREMIUM: 2,
}

// RequireSubscription aborts with 402 if the authenticated user's subscription
// level is below minLevel. Admins always pass regardless of their tier.
func RequireSubscription(minLevel string) gin.HandlerFunc {
	required := subscriptionOrder[minLevel]
	return func(c *gin.Context) {
		// Admins bypass subscription gates.
		if role, _ := c.Get("user_role"); role == config.ADMIN_ROLE {
			c.Next()
			return
		}

		sub, _ := c.Get("user_subscription")
		subStr, _ := sub.(string)
		if subscriptionOrder[subStr] < required {
			c.JSON(402, gin.H{
				"success":  false,
				"error":    "subscription_required",
				"message":  "This feature requires a " + minLevel + " subscription. Upgrade to unlock it.",
				"required": minLevel,
				"current":  subStr,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
