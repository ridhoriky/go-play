package transaction

import (
	"context"
	"ne-project/internal/dto"
	"ne-project/internal/models"

	"github.com/jmoiron/sqlx"
)

type TransactionRepositoryItf interface {
	Checkout(
		ctx context.Context,
		req *dto.CheckoutRequest,
	) (*models.Transaction, error)
	GetToday(ctx context.Context) (*dto.TodaySummaryResponse, error)
}
type TransactionRepository struct {
	db  *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepositoryItf {
	return &TransactionRepository{
		db: db,
	}
}