package transaction

import (
	"context"
	"ne-project/internal/dto"
	"ne-project/internal/repositories/transaction"
)

type TransactionServiceItf interface {
	Checkout(
		ctx context.Context,
		req *dto.CheckoutRequest,
	) (*dto.TransactionResponse, error)
	GetToday(ctx context.Context) ([]dto.TransactionResponse, error)
}

	
type transactionService struct {
	transactionRepository transaction.TransactionRepositoryItf

}

func NewTransactionService(transactionRepository transaction.TransactionRepositoryItf) TransactionServiceItf {
	return &transactionService{
		transactionRepository: transactionRepository,
	}
}
