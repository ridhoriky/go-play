package transaction

import (
	"context"
	"database/sql"
	"ne-project/internal/dto"
	"ne-project/internal/models"
)

type TransactionRepositoryItf interface {
	Checkout(
		ctx context.Context,
		req *dto.CheckoutRequest,
	) (*models.Transaction, error)
	GetToday(ctx context.Context) ([]models.Transaction, error)
}
type TransactionRepository struct {
	db  *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepositoryItf {
	return &TransactionRepository{
		db: db,
	}
}