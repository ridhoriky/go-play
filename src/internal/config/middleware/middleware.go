package middleware

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var (
	onceMiddleware = &sync.Once{}
	middlewareInst Middleware
)

// Middleware defines the middleware interface
type Middleware interface {
	Handler() gin.HandlerFunc
	CORS() gin.HandlerFunc
	Logger() gin.HandlerFunc
}

type middleware struct {
	log zerolog.Logger
}

// RateLimiterOptions holds rate limiter configuration
type RateLimiterOptions struct {
	Command string `yaml:"command"`
	Limit   int    `yaml:"limit"`
}

// InitMiddleware initializes the middleware
func InitMiddleware(log zerolog.Logger) Middleware {
	onceMiddleware.Do(func() {
		middlewareInst = &middleware{
			log: log,
		}
	})
	return middlewareInst
}

// Handler returns the main handler
func (m *middleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// CORS returns the CORS handler
func (m *middleware) CORS() gin.HandlerFunc {
	allowedOrigins := getAllowedOrigins()

	return func(c *gin.Context) {
		origin := strings.TrimSpace(c.GetHeader("Origin"))
		if origin == "" || !isOriginAllowed(origin, allowedOrigins) {
			c.Next()
			return
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Vary", "Origin")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders(c.Request.Header.Get("Access-Control-Request-Headers")))
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (m *middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		logger := m.log.Info()

		switch {
		case status >= http.StatusInternalServerError:
			logger = m.log.Error()
		case status >= http.StatusBadRequest:
			logger = m.log.Warn()
		}

		event := logger.
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("request_uri", c.Request.RequestURI).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Int("status", status).
			Dur("duration_ms", duration)

		if query := c.Request.URL.RawQuery; query != "" {
			event = event.Str("query", query)
		}

		if c.Errors.Last() != nil {
			event = event.Str("errors", c.Errors.String())
		}

		event.Msg("http request")
	}
}

func allowedHeaders(requestHeaders string) string {
	requestHeaders = strings.TrimSpace(requestHeaders)
	if requestHeaders != "" {
		return requestHeaders
	}

	return "Content-Type, Authorization, Accept, Origin, User-Agent, X-Requested-With, X-CSRF-Token, Cache-Control"
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		allowed = strings.TrimSpace(allowed)
		if allowed == "*" || strings.EqualFold(allowed, origin) {
			return true
		}
	}

	return false
}

func getAllowedOrigins() []string {
	allowedOriginEnv := strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS"))
	if allowedOriginEnv == "" {
		return []string{"http://localhost:3000", "http://localhost:8080"}
	}

	origins := strings.Split(allowedOriginEnv, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return origins
}
