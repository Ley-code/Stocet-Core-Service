package utils

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"stocet-core-service/internal/config"

	"github.com/gin-gonic/gin"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// SuccessResponse sends a successful JSON response
func SuccessResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// ErrorResponse sends an error JSON response
func ErrorResponse(c *gin.Context, statusCode int, err error, message string) {
	var errString string
	if err != nil {
		if config.AppConfig != nil && config.AppConfig.Environment == "production" {
			log.Printf("[error] status=%d message=%q err=%v", statusCode, message, err)
			errString = "internal_error"
		} else {
			errString = err.Error()
		}
	}
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   errString,
		Message: message,
	})
}

// GetUserID returns the authenticated user ID from Gin context.
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	uid, ok := userID.(string)
	return uid, ok
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, nil, message)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, nil, message)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, nil, message)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, nil, message)
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *gin.Context, err error, message string) {
	ErrorResponse(c, http.StatusInternalServerError, err, message)
}

// PaginationParams represents pagination parameters
type PaginationParams struct {
	Page     int
	PageSize int
	Skip     int
	Limit    int
}

// ParsePagination parses pagination parameters from query
func ParsePagination(c *gin.Context, defaultPageSize, maxPageSize int) PaginationParams {
	page := 1
	pageSize := defaultPageSize

	if p := c.Query("page"); p != "" {
		if parsedPage := parseInt(p); parsedPage > 0 {
			page = parsedPage
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsedPageSize := parseInt(ps); parsedPageSize > 0 {
			if parsedPageSize > maxPageSize {
				pageSize = maxPageSize
			} else {
				pageSize = parsedPageSize
			}
		}
	}

	skip := (page - 1) * pageSize
	limit := pageSize

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Skip:     skip,
		Limit:    limit,
	}
}

// parseInt is a helper to parse integers (simplified version)
func parseInt(s string) int {
	var result int
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		} else {
			return 0
		}
	}
	return result
}

// ParseInt is a public version of parseInt
func ParseInt(s string) int {
	return parseInt(s)
}

// ParseDate parses a date string in multiple formats
// Supports: RFC3339, "2006-01-02", "2006-01-02T15:04:05Z", "2006-01-02 15:04:05"
func ParseDate(dateStr string) (time.Time, error) {
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("empty date string")
	}

	// Try RFC3339 format first (full ISO 8601 with timezone)
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t, nil
	}

	// Try RFC3339 without timezone
	if t, err := time.Parse("2006-01-02T15:04:05", dateStr); err == nil {
		return t, nil
	}

	// Try simple date format "2006-01-02" (set to start of day UTC)
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
	}

	// Try date with time "2006-01-02 15:04:05"
	if t, err := time.Parse("2006-01-02 15:04:05", dateStr); err == nil {
		return t, nil
	}

	// Try date with time and timezone "2006-01-02 15:04:05 -0700"
	if t, err := time.Parse("2006-01-02 15:04:05 -0700", dateStr); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
