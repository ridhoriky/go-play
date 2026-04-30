package middleware

import (
	"ne-project/src/internal/config/token"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Middleware defines the middleware interface

type Middleware struct {
	log      zerolog.Logger
	tokenSvc *token.Token
	limit    int
	window   time.Duration
}

// InitMiddleware initializes the middleware
func InitMiddleware(log zerolog.Logger, tokenSvc *token.Token, limit int, windows time.Duration) *Middleware {
	return &Middleware{
		log:      log,
		tokenSvc: tokenSvc,
		limit:    limit,
		window:   windows,
	}

}

func (m *Middleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
