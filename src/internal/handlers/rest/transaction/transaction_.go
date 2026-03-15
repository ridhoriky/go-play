package transaction

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"

	"github.com/gin-gonic/gin"
)

func (h *transactionHandler) GetToday(c *gin.Context) {

	ctx := c.Request.Context()

	result, err := h.transactionService.GetToday(ctx)
	if err != nil {
		helpers.ResponseError(
			c.Writer,
			&dto.Error{
				Code:    http.StatusBadRequest,
				Message: preference.ErrNoProductCreated,
			},
		)
		return
	}

	helpers.ResponseSuccess(
		c.Writer,
		http.StatusOK,
		"success",
		result,
	)
}
