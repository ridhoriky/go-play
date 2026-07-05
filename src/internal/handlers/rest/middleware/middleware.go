package middleware

import (
	"ne-project/src/internal/config/token"

	"github.com/rs/zerolog"
)

// Middleware defines the middleware interface

type Middleware struct {
	log       zerolog.Logger
	tokenSvc  token.TokenServiceItf
	rateLimit *RateLimiterOptions
}

// InitMiddleware initializes the middleware
func InitMiddleware(log *zerolog.Logger, tokenSvc token.TokenServiceItf, rateLimit *RateLimiterOptions) *Middleware {
	return &Middleware{
		log:       *log,
		tokenSvc:  tokenSvc,
		rateLimit: rateLimit,
	}

}
