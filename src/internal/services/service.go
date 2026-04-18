package services

import (
	"ne-project/src/internal/repositories"
	"ne-project/src/internal/services/category"
	"ne-project/src/internal/services/product"
	"ne-project/src/internal/services/report"
	"ne-project/src/internal/services/transaction"
)

type Services struct {
	Category    category.CategoryServiceItf
	Product     product.ProductServiceItf
	Transaction transaction.TransactionServiceItf
	Report      report.ReportServiceItf
}

func NewServices(repositories *repositories.Repositories) *Services {
	return &Services{
		Category:    category.NewCategoryService(repositories.Category),
		Product:     product.NewProductService(repositories.Product),
		Transaction: transaction.NewTransactionService(repositories.Transaction),
		Report:      report.NewReportService(repositories.Report),
	}
}
