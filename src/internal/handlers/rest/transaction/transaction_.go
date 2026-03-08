package transaction

import (
	"net/http"

	"ne-project/src/internal/models/dto"

	"github.com/gin-gonic/gin"
)

func (h *transactionHandler) GetToday(c *gin.Context) {

	ctx := c.Request.Context()

	result, err := h.transactionService.GetToday(ctx)
	if err != nil {
		dto.ResponseError(
			c.Writer,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	dto.ResponseSuccess(
		c.Writer,
		http.StatusOK,
		"success",
		result,
	)
}
