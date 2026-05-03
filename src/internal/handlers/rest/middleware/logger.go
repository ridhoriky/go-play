package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (m *Middleware) SetComponent(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := zerolog.Ctx(c.Request.Context()).With().Str("component", name).Logger()
		c.Request = c.Request.WithContext(l.WithContext(c.Request.Context()))
		c.Next()
	}
}

func (m *Middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header("X-Request-ID", requestID)

		reqLogger := m.log.With().
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("client_ip", c.ClientIP()).
			Logger()

		ctx := reqLogger.WithContext(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()

		event := reqLogger.Info()
		if status >= http.StatusInternalServerError {
			event = reqLogger.Error()
		} else if status >= http.StatusBadRequest {
			event = reqLogger.Warn()
		}

		event = event.
			Int("status", status).
			Int64("duration_ms", duration.Milliseconds()).
			Str("user_agent", c.Request.UserAgent())

		if query := c.Request.URL.RawQuery; query != "" {
			event = event.Str("query", query)
		}

		if len(c.Errors) > 0 {
			event = event.Str("errors", c.Errors.String())
		}

		event.Msg("http request")
	}
}
