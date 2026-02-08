package transaction

import (
	"ne-project/internal/services/transaction"

	"github.com/gin-gonic/gin"
)


type TransactionHandlerItf interface {
	RegisterRoutes(r *gin.Engine)
}

type transactionHandler struct {
	transactionService transaction.TransactionServiceItf
	
}

func NewTransactionHandler(transactionService transaction.TransactionServiceItf) TransactionHandlerItf {
	return &transactionHandler{
		transactionService: transactionService,
	}
}

func (h *transactionHandler) RegisterRoutes(r *gin.Engine) {
	transactionRoutes := r.Group("/transactions")
	{
		transactionRoutes.POST("/checkout", h.Checkout)
		transactionRoutes.GET("/today", h.GetToday)
	}
}