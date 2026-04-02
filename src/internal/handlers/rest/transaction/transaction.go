package transaction

import (
	"ne-project/src/internal/services/transaction"

	"github.com/gin-gonic/gin"
)

type TransactionHandlerItf interface {
	RegisterRoutes(r *gin.RouterGroup)
}

type transactionHandler struct {
	transactionService transaction.TransactionServiceItf
}

func NewTransactionHandler(transactionService transaction.TransactionServiceItf) TransactionHandlerItf {
	return &transactionHandler{
		transactionService: transactionService,
	}
}

func (h *transactionHandler) RegisterRoutes(r *gin.RouterGroup) {
	transactionRoutes := r.Group("/transactions")
	{
		transactionRoutes.POST("/", h.Checkout)
		transactionRoutes.GET("/:id", h.GetByID)
	}
}
