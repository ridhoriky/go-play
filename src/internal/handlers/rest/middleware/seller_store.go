package middleware

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) RequireSellerStore() gin.HandlerFunc {
	return func(c *gin.Context) {
		if storeID := c.GetString("store_id"); storeID != "" {
			c.Next()
			return
		}

		userID := c.GetString("user_id")
		if userID == "" {
			helpers.ResponseError(c, dto.NewError(http.StatusUnauthorized, preference.ErrUnauthorized))
			c.Abort()
			return
		}

		if m.storeSvc != nil {
			storeRes, err := m.storeSvc.GetMyStore(c.Request.Context(), userID)
			if err == nil && storeRes != nil {
				c.Set("store_id", storeRes.ID)
				c.Next()
				return
			}
		}

		helpers.ResponseError(c, dto.NewError(http.StatusForbidden, "Seller store not found"))
		c.Abort()
	}
}
