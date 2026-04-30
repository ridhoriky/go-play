package middleware

import (
	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) CORS() gin.HandlerFunc {
	allowedOrigins := getAllowedOrigins()

	return func(c *gin.Context) {
		origin := strings.TrimSpace(c.GetHeader("Origin"))

		if origin != "" && !isOriginAllowed(origin, allowedOrigins) {
			err := &dto.Error{
				Code:    http.StatusForbidden,
				Message: "Origin not allowed",
			}
			helpers.ResponseError(c.Writer, err)
			c.Abort()
			return
		}

		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", allowedHeaders(c.Request.Header.Get("Access-Control-Request-Headers")))
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Security headers
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self';")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
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
