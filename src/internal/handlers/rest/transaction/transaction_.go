package transaction

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"

	"github.com/gin-gonic/gin"
)

func (h *transactionHandler) Checkout(c *gin.Context) {

	ctx := c.Request.Context()

	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrInvalidReqBody,
		})
		return
	}

	result, err := h.transactionService.Checkout(ctx, &req)
	if err != nil {
		helpers.ResponseError(c.Writer, err)
		return
	}

	helpers.ResponseSuccess(c.Writer, http.StatusOK, "Checkout successfully", result)

}
