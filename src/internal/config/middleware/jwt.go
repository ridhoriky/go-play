package middleware

import (
	"net/http"
	"strings"

	"ne-project/src/internal/config/token"
	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"

	"github.com/gin-gonic/gin"
)

func (m *middleware) JWTAuth(tokenSvc token.Token) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			helpers.ResponseError(c.Writer, &dto.Error{
				Code:    http.StatusUnauthorized,
				Message: preference.ErrMissingAuthHeader,
			})
			c.Abort()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			helpers.ResponseError(c.Writer, &dto.Error{
				Code:    http.StatusUnauthorized,
				Message: preference.ErrInvalidAuthFormat,
			})
			c.Abort()
			return
		}

		claims, err := tokenSvc.ValidateAccessToken(parts[1])
		if err != nil {
			helpers.ResponseError(c.Writer, &dto.Error{
				Code:    http.StatusUnauthorized,
				Message: preference.ErrInvalidCredentials,
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("user_name", claims.Name)
		c.Next()
	}
}
