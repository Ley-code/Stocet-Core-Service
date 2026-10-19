package middleware

import (
	"strings"

	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// OptionalAuthMiddleware lets guests through. A valid Bearer token still populates the
// user context keys (so personalised data keeps working); a missing or invalid one is ignored.
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		parts := strings.Split(c.GetHeader("Authorization"), " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			if claims, err := utils.ValidateToken(parts[1]); err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
				c.Set("user_role", claims.Role)
				c.Set("user_subscription", claims.SubscriptionStatus)
				if claims.OrgID != "" {
					c.Set("user_org_id", claims.OrgID)
				}
			}
		}
		c.Next()
	}
}
