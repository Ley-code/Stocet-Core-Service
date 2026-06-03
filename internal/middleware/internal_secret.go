package middleware

import (
	"crypto/subtle"
	"stocet-core-service/internal/config"
	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
)

// InternalSecretMiddleware protects internal service-to-service endpoints.
func InternalSecretMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secret := config.AppConfig.InternalSecret
		if secret == "" {
			utils.InternalServerError(c, nil, "Internal endpoint is not configured")
			c.Abort()
			return
		}

		provided := c.GetHeader("X-Internal-Secret")
		if provided == "" {
			auth := c.GetHeader("Authorization")
			if len(auth) > 7 && auth[:7] == "Bearer " {
				provided = auth[7:]
			}
		}

		if subtle.ConstantTimeCompare([]byte(provided), []byte(secret)) != 1 {
			utils.Unauthorized(c, "Invalid internal credentials")
			c.Abort()
			return
		}

		c.Next()
	}
}
