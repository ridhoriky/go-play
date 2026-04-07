package repositories

import (
	"ne-project/src/internal/repositories/category"
	"ne-project/src/internal/repositories/product"
	"ne-project/src/internal/repositories/report"
	"ne-project/src/internal/repositories/transaction"

	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	Category    category.CategoryRepositoryItf
	Product     product.ProductRepositoryItf
	Transaction transaction.TransactionRepositoryItf
	Report      report.ReportRepositoryItf
}

func NewRepository(db *sqlx.DB) *Repositories {
	return &Repositories{
		Category:    category.NewCategoryRepository(db),
		Product:     product.NewProductRepository(db),
		Transaction: transaction.NewTransactionRepository(db),
		Report:      report.NewReportRepository(db),
	}
}
