package transaction

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// CreateTransaction godoc
// @Summary      Create Transaction
// @Description  Create a transaction with multiple items
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param 		 transaction 	body 	dto.CreateTransactionRequest true "Transaction data"
// @Success 	 201 		{object} 	dto.APIResponse{data=dto.TransactionDetailResponse}
// @Failure      400  		{object} 	dto.APIResponse
// @Failure      404  		{object}  	dto.APIResponse
// @Router       /transactions [post]
func (h *transactionHandler) Checkout(c *gin.Context) {

	ctx := c.Request.Context()

	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg(preference.ErrInvalidReqBody)
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
