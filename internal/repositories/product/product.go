package product

import (
	"context"
	"database/sql"
	"ne-project/internal/dto"
	"ne-project/internal/models"
)

type ProductRepositoryItf interface {
	GetAll(ctx context.Context) ([]dto.ProductResponse, error)
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id int) (*dto.ProductResponse, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
}

type ProductRepository struct {
	db  *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepositoryItf {
	return &ProductRepository{
		db: db,
	}
}