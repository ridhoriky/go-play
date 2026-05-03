package middleware

import (
	"errors"
	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var appErr *dto.Error

		if !errors.As(err, &appErr) {
			appErr = helpers.ParsePgError(err)
		}

		if c.Writer.Written() {
			return
		}

		resp := dto.APIResponse{
			Status:  appErr.Code,
			Message: appErr.Message,
			Error:   preference.ErrorCodeByHTTPStatus[appErr.Code],
		}

		c.AbortWithStatusJSON(appErr.Code, resp)
	}
}
