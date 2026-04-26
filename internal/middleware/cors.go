package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles CORS headers for browser clients (SPA on another origin).
// Note: Browsers forbid Access-Control-Allow-Origin: * together with Allow-Credentials: true.
// This API uses Bearer tokens (not cookie auth), so credentials are omitted and * is valid.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
