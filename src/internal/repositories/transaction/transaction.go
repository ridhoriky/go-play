package transaction

import (
	"context"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"

	"github.com/jmoiron/sqlx"
)

type TransactionRepositoryItf interface {
	Checkout(ctx context.Context, req *dto.CreateTransactionRequest) (*entity.TransactionWithDetails, error)
	GetTransactionByID(ctx context.Context, id string) (*entity.TransactionWithDetails, error)
	UpdateStatus(ctx context.Context, id string, status entity.TransactionStatus, items []entity.TransactionDetail) error
}
type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepositoryItf {
	return &TransactionRepository{
		db: db,
	}
}
