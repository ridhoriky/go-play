package transaction

import (
	"context"
	"ne-project/src/internal/models/dto"

	"github.com/jmoiron/sqlx"
)

type TransactionRepositoryItf interface {
	GetToday(ctx context.Context) (*dto.TransactionListResponse, error)
}
type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepositoryItf {
	return &TransactionRepository{
		db: db,
	}
}
