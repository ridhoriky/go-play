package product

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"

	"github.com/jmoiron/sqlx"
)

type ProductRepositoryItf interface {
	CreateMultiple(ctx context.Context, products []entity.Product) ([]entity.Product, error)
	GetAll(ctx context.Context, req *dto.GetProductsQuery) ([]entity.ProductWithCategory, int, error)
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id string) (*entity.Product, string, error)
	Update(ctx context.Context, id string, product *entity.Product) error
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
