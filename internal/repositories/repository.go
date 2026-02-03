package repositories

import (
	"database/sql"
	"ne-project/internal/repositories/category"
	"ne-project/internal/repositories/product"
)

type Repositories struct {
	Category category.CategoryRepositoryItf
	Product  product.ProductRepositoryItf
}

func NewRepository(db *sql.DB) *Repositories {
	return &Repositories{
		Category: category.NewCategoryRepository(db),
		Product:  product.NewProductRepository(db),
	}
}