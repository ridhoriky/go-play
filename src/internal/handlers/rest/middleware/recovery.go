package middleware

import (
	"net/http"
	"runtime/debug"

	"ne-project/src/internal/preference"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				m.log.Error().
					Interface("error", err).
					Str("stack", string(stack)).
					Msg("panic recovered")

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"status":  http.StatusInternalServerError,
					"message": "Internal Server Error",
					"error":   preference.ErrorCodeByHTTPStatus[http.StatusInternalServerError],
				})
			}
		}()
		c.Next()
	}
}
