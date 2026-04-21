package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (m *middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := m.log.WithContext(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)
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
