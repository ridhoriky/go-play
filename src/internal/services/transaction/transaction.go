package transaction

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/repositories/transaction"
)

type TransactionServiceItf interface {
	GetToday(ctx context.Context) (*dto.TransactionListResponse, error)
}

type transactionService struct {
	transactionRepository transaction.TransactionRepositoryItf
}

func NewTransactionService(transactionRepository transaction.TransactionRepositoryItf) TransactionServiceItf {
	return &transactionService{
		transactionRepository: transactionRepository,
	}
}
