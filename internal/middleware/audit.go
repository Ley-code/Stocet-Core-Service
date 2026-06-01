package middleware

import (
	"net/http"

	"stocet-core-service/internal/models"
	"stocet-core-service/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// auditResponseWriter wraps gin.ResponseWriter to capture the response status code.
type auditResponseWriter struct {
	gin.ResponseWriter
	statusCode int
}

func (w *auditResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *auditResponseWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

// AuditMiddleware logs every successful (2xx/3xx) admin mutation asynchronously.
// GET/HEAD/OPTIONS requests are never logged.
func AuditMiddleware(repo *repositories.AuditLogRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
			c.Next()
			return
		}

		arw := &auditResponseWriter{ResponseWriter: c.Writer, statusCode: http.StatusOK}
		c.Writer = arw

		c.Next()

		// Only log successful mutations
		if arw.statusCode >= 400 {
			return
		}

		actorIDRaw, _ := c.Get("user_id")
		actorEmail, _ := c.Get("user_email")
		actorIDStr, _ := actorIDRaw.(string)
		actorUID, err := uuid.Parse(actorIDStr)
		if err != nil {
			return
		}

		entry := &models.AuditLog{
			ActorID:    actorUID,
			ActorEmail: func() string { s, _ := actorEmail.(string); return s }(),
			Method:     method,
			Path:       c.FullPath(),
			ResourceID: c.Param("id"),
			StatusCode: arw.statusCode,
			IPAddress:  c.ClientIP(),
		}

		go func() { _ = repo.Create(entry) }()
	}
}
