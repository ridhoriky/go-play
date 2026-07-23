package middleware

import (
	"ne-project/src/internal/config/token"
	"ne-project/src/internal/services/store"

	"github.com/rs/zerolog"
)

// Middleware defines the middleware interface

type Middleware struct {
	log       zerolog.Logger
	tokenSvc  token.TokenServiceItf
	rateLimit *RateLimiterOptions
	storeSvc  store.StoreServiceItf
}

// InitMiddleware initializes the middleware
func InitMiddleware(log *zerolog.Logger, tokenSvc token.TokenServiceItf, rateLimit *RateLimiterOptions, storeSvc store.StoreServiceItf) *Middleware {
	return &Middleware{
		log:       *log,
		tokenSvc:  tokenSvc,
		rateLimit: rateLimit,
		storeSvc:  storeSvc,
	}
}
