package middleware

import (
	"net/http"
	"sync"
	"time"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	mu          sync.Mutex
	count       int
	windowStart time.Time
}

var rateLimits sync.Map

func RateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if limit <= 0 || window <= 0 {
			c.Next()
			return
		}

		key := c.ClientIP() + ":" + c.FullPath()
		value, _ := rateLimits.LoadOrStore(key, &rateLimitEntry{
			count:       0,
			windowStart: time.Now(),
		})

		entry := value.(*rateLimitEntry)
		entry.mu.Lock()
		defer entry.mu.Unlock()

		now := time.Now()
		if now.Sub(entry.windowStart) > window {
			entry.count = 0
			entry.windowStart = now
		}

		entry.count++
		if entry.count > limit {
			helpers.ResponseError(c.Writer, &dto.Error{
				Code:    http.StatusTooManyRequests,
				Message: "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
