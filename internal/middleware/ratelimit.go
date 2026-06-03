package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"stocet-core-service/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	limiterMu sync.Mutex
	limiters  = make(map[string]*ipLimiter)
)

func getIPLimiter(key string, r rate.Limit, burst int) *rate.Limiter {
	limiterMu.Lock()
	defer limiterMu.Unlock()

	entry, ok := limiters[key]
	if !ok {
		entry = &ipLimiter{
			limiter:  rate.NewLimiter(r, burst),
			lastSeen: time.Now(),
		}
		limiters[key] = entry
	}
	entry.lastSeen = time.Now()
	return entry.limiter
}

func cleanupLimiters() {
	limiterMu.Lock()
	defer limiterMu.Unlock()
	cutoff := time.Now().Add(-10 * time.Minute)
	for key, entry := range limiters {
		if entry.lastSeen.Before(cutoff) {
			delete(limiters, key)
		}
	}
}

func RateLimitByIP(requestsPerMinute int, burst int) gin.HandlerFunc {
	if burst <= 0 {
		burst = requestsPerMinute
	}
	r := rate.Every(time.Minute / time.Duration(requestsPerMinute))

	return func(c *gin.Context) {
		key := c.ClientIP() + ":" + c.FullPath()
		limiter := getIPLimiter(key, r, burst)
		if !limiter.Allow() {
			utils.ErrorResponse(c, http.StatusTooManyRequests, nil, "Too many requests")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RateLimitWhenQuery applies IP rate limiting only when the given query param is non-empty.
func RateLimitWhenQuery(param string, requestsPerMinute, burst int) gin.HandlerFunc {
	limiter := RateLimitByIP(requestsPerMinute, burst)
	return func(c *gin.Context) {
		if strings.TrimSpace(c.Query(param)) == "" {
			c.Next()
			return
		}
		limiter(c)
	}
}

func init() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			cleanupLimiters()
		}
	}()
}
