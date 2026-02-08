package services

import (
	"ne-project/internal/repositories"
	"ne-project/internal/services/category"
	"ne-project/internal/services/product"
	"ne-project/internal/services/transaction"
)

type Services struct {
	Category category.CategoryServiceItf
	Product product.ProductServiceItf
	Transaction transaction.TransactionServiceItf
}

func NewServices(repositories *repositories.Repositories) *Services {
	return &Services{
		Category: category.NewCategoryService(
			repositories.Category,
		),
		Product: product.NewProductService(
			repositories.Product,
		),
		Transaction: transaction.NewTransactionService(
			repositories.Transaction,
		),
	}
}