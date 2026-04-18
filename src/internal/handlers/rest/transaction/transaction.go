package transaction

import (
	"ne-project/src/internal/services/transaction"

	"github.com/gin-gonic/gin"
)

type TransactionHandlerItf interface {
	Checkout(c *gin.Context)
	GetByID(c *gin.Context)
	UpdateTransactionStatus(c *gin.Context)
}

type transactionHandler struct {
	transactionService transaction.TransactionServiceItf
}

func NewTransactionHandler(transactionService transaction.TransactionServiceItf) TransactionHandlerItf {
	return &transactionHandler{
		transactionService: transactionService,
	}
}
