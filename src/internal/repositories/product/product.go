package product

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"

	"github.com/jmoiron/sqlx"
)

type ProductRepositoryItf interface {
	CreateMultiple(ctx context.Context, products []entity.Product) ([]dto.ProductDTO, error)
	GetAll(ctx context.Context, req dto.ProductFilterRequest) ([]dto.ProductResponse, error)
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id string) (*dto.ProductResponse, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id string) error
}

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepositoryItf {
	return &ProductRepository{
		db: db,
	}
}
