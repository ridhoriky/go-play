package rest

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/transaction"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionService transaction.TransactionServiceItf
}

func NewTransactionHandler(transactionService transaction.TransactionServiceItf) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

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
func (h *TransactionHandler) Checkout(c *gin.Context) {

	ctx := c.Request.Context()

	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	result, err := h.transactionService.Checkout(ctx, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Checkout successfully", result)

}

// GetTransactionByID godoc
// @Summary      Get Transaction
// @Description  Get transaction by their id
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param 		 id   path 	string true "Category ID (UUID)" format(uuid)
// @Success      200  {object}  dto.APIResponse
// @Failure      400  {object} 	dto.APIResponse
// @Failure      404  {object}  dto.APIResponse
// @Router       /transactions/{id} [get]
func (h *TransactionHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidTransactionID))
		return
	}
	transaction, err := h.transactionService.GetTransactionByID(ctx, id.String())
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", transaction)
}

// UpdateTransactionStatus godoc
// @Summary      Update Transaction Status
// @Description  Update transaction Status by their id
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param 		 id   path 	string true "Transaction ID (UUID)" format(uuid)
// @Success      200  {object}  dto.APIResponse
// @Failure      400  {object} 	dto.APIResponse
// @Failure      404  {object}  dto.APIResponse
// @Router       /transactions/{id}/status [patch]
func (h *TransactionHandler) UpdateTransactionStatus(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidTransactionID))
		return
	}

	var status *dto.UpdateTransactionStatusRequest
	if err = c.ShouldBindJSON(&status); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	transaction, err := h.transactionService.UpdateStatus(ctx, id.String(), status)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", transaction)
}
