package transaction

import (
	"context"

	"ne-project/internal/models/dto"
	"ne-project/internal/repositories/transaction"
)

type TransactionServiceItf interface {
	GetToday(ctx context.Context) (*dto.TodaySummaryResponse, error)
}

type transactionService struct {
	transactionRepository transaction.TransactionRepositoryItf
}

func NewTransactionService(transactionRepository transaction.TransactionRepositoryItf) TransactionServiceItf {
	return &transactionService{
		transactionRepository: transactionRepository,
	}
}
