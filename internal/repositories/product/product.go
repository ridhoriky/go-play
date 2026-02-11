package product

import (
	"context"
	"ne-project/internal/dto"
	"ne-project/internal/models"

	"github.com/jmoiron/sqlx"
)

type ProductRepositoryItf interface {
	CreateMultiple(ctx context.Context, products []models.Product) ([]dto.ProductDTO, error)
	GetAll(ctx context.Context, req dto.ProductFilterRequest) ([]dto.ProductResponse, error)
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id int) (*dto.ProductResponse, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
}

type ProductRepository struct {
	db  *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepositoryItf {
	return &ProductRepository{
		db: db,
	}
}