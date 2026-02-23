package transaction

import (
	"context"
	"ne-project/internal/models/dto"

	"github.com/jmoiron/sqlx"
)

type TransactionRepositoryItf interface {
	GetToday(ctx context.Context) (*dto.TodaySummaryResponse, error)
}
type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepositoryItf {
	return &TransactionRepository{
		db: db,
	}
}
