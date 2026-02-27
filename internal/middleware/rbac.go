package middleware

import (
	"stocet-core-service/internal/config"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// RequireRole checks if the user has the required role
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			utils.Unauthorized(c, "User role not found")
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok || role != requiredRole {
			utils.Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin is a convenience middleware for admin-only routes
func RequireAdmin() gin.HandlerFunc {
	return RequireRole(config.ADMIN_ROLE)
}

// RequireAnyRole checks if the user has any of the required roles
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			utils.Unauthorized(c, "User role not found")
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			utils.Forbidden(c, "Invalid user role")
			c.Abort()
			return
		}

		for _, requiredRole := range roles {
			if role == requiredRole {
				c.Next()
				return
			}
		}

		utils.Forbidden(c, "Insufficient permissions")
		c.Abort()
	}
}
