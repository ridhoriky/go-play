package transaction

import (
	"encoding/json"
	"ne-project/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *transactionHandler) Checkout(c *gin.Context) {
	ctx := c.Request.Context()
	var t dto.CheckoutRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&t); err != nil {
		dto.ResponseError(c.Writer, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := t.Validate(); err != nil {
		dto.ResponseError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}
	responses, err := h.transactionService.Checkout(ctx, &t); 
	if err != nil {
		dto.ResponseError(c.Writer, http.StatusInternalServerError, "Failed to create transaction")
		return
	}
	dto.ResponseSuccess(c.Writer, http.StatusCreated, "Transaction created successfully", responses)
}

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

