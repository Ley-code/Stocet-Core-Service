package middleware

import (
	"stocet-core-service/internal/config"
	"stocet-core-service/internal/repositories"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

var adminUserRepo *repositories.UserRepository

// SetAdminUserRepository enables DB-backed admin checks in RequireAdmin.
func SetAdminUserRepository(repo *repositories.UserRepository) {
	adminUserRepo = repo
}

func userHasRole(c *gin.Context, requiredRole string) bool {
	if adminUserRepo != nil {
		userID, ok := utils.GetUserID(c)
		if !ok {
			return false
		}
		user, err := adminUserRepo.GetByID(userID)
		if err != nil {
			return false
		}
		c.Set("user_role", user.Role)
		c.Set("user_subscription", user.SubscriptionStatus)
		return string(user.Role) == requiredRole
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		return false
	}
	role, ok := userRole.(string)
	return ok && role == requiredRole
}

// RequireRole checks if the user has the required role
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !userHasRole(c, requiredRole) {
			if _, ok := utils.GetUserID(c); !ok {
				utils.Unauthorized(c, "User role not found")
			} else {
				utils.Forbidden(c, "Insufficient permissions")
			}
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
		for _, requiredRole := range roles {
			if userHasRole(c, requiredRole) {
				c.Next()
				return
			}
		}

		if _, ok := utils.GetUserID(c); !ok {
			utils.Unauthorized(c, "User role not found")
		} else {
			utils.Forbidden(c, "Insufficient permissions")
		}
		c.Abort()
	}
}
