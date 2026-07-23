package middleware

import (
	"net/http"
	"strings"

	"ne-project/src/internal/config/token"
	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func (m *Middleware) JWTAuth(tokenSvc token.TokenServiceItf) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			helpers.ResponseError(c, dto.NewError(http.StatusUnauthorized, preference.ErrMissingAuthHeader))
			c.Abort()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			helpers.ResponseError(c, dto.NewError(http.StatusUnauthorized, preference.ErrInvalidAuthFormat))
			c.Abort()
			return
		}

		claims, err := tokenSvc.ValidateAccessToken(parts[1])
		if err != nil {
			helpers.ResponseError(c, dto.NewError(http.StatusUnauthorized, preference.ErrInvalidCredentials))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("user_name", claims.Name)
		c.Set("store_id", claims.StoreID)

		ctx := c.Request.Context()
		logger := zerolog.Ctx(ctx).With().
			Str("user_id", claims.UserID).
			Str("user_role", claims.Role).
			Logger()
		c.Request = c.Request.WithContext(logger.WithContext(ctx))

		c.Next()
	}
}

func (m *Middleware) OptionalJWTAuth(tokenSvc token.TokenServiceItf) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			c.Next()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.Next()
			return
		}

		claims, err := tokenSvc.ValidateAccessToken(parts[1])
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("user_name", claims.Name)
		c.Set("store_id", claims.StoreID)

		// Enrich context logger with user data
		ctx := c.Request.Context()
		logger := zerolog.Ctx(ctx).With().
			Str("user_id", claims.UserID).
			Str("user_role", claims.Role).
			Logger()
		c.Request = c.Request.WithContext(logger.WithContext(ctx))

		c.Next()
	}
}
