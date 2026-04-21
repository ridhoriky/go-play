package middleware

import (
	"ne-project/src/internal/config/token"
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
	JWTAuth(tokenSvc token.Token) gin.HandlerFunc
	RateLimiter(limit int, window time.Duration) gin.HandlerFunc
}

type middleware struct {
	log      zerolog.Logger
	tokenSvc token.Token
	limit    int
	window   time.Duration
}

// InitMiddleware initializes the middleware
func InitMiddleware(log zerolog.Logger, tokenSvc token.Token, limit int, windows time.Duration) Middleware {
	onceMiddleware.Do(func() {
		middlewareInst = &middleware{
			log:      log,
			tokenSvc: tokenSvc,
			limit:    limit,
			window:   windows,
		}
	})
	return middlewareInst
}

func (m *middleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
