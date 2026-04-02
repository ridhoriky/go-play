package transaction

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/repositories/transaction"
)

type TransactionServiceItf interface {
	Checkout(ctx context.Context, req *dto.CreateTransactionRequest) (*dto.TransactionDetailResponse, error)
	GetTransactionByID(ctx context.Context, id string) (*dto.TransactionDetailResponse, error)
}

type transactionService struct {
	transactionRepository transaction.TransactionRepositoryItf
}

func NewTransactionService(transactionRepository transaction.TransactionRepositoryItf) TransactionServiceItf {
	return &transactionService{
		transactionRepository: transactionRepository,
	}
}
