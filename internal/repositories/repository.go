package repositories

import (
	"database/sql"
	"ne-project/internal/repositories/category"
	"ne-project/internal/repositories/product"
	"ne-project/internal/repositories/transaction"
)

type Repositories struct {
	Category category.CategoryRepositoryItf
	Product  product.ProductRepositoryItf
	Transaction  transaction.TransactionRepositoryItf
}

func NewRepository(db *sql.DB) *Repositories {
	return &Repositories{
		Category: category.NewCategoryRepository(db),
		Product:  product.NewProductRepository(db),
		Transaction: transaction.NewTransactionRepository(db),
	}
}