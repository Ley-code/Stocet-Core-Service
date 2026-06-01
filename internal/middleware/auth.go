package middleware

import (
	"strings"

	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens and populates user context keys.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("user_subscription", claims.SubscriptionStatus)
		if claims.OrgID != "" {
			c.Set("user_org_id", claims.OrgID)
		}

		c.Next()
	}
}
