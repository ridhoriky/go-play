package middleware

import (
	"net/http"
	"slices"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"

	"github.com/gin-gonic/gin"
)

// RequireRole checks if the user has one of the required roles
func (m *Middleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			helpers.ResponseError(c, dto.NewError(http.StatusUnauthorized, preference.ErrUnauthorized))
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			helpers.ResponseError(c, dto.NewError(http.StatusForbidden, preference.ErrForbidden))
			c.Abort()
			return
		}

		hasRole := slices.Contains(roles, roleStr)

		if !hasRole {
			helpers.ResponseError(c, dto.NewError(http.StatusForbidden, preference.ErrForbidden))
			c.Abort()
			return
		}

		c.Next()
	}
}
