package services

import (
	"ne-project/internal/repositories"
	"ne-project/internal/services/category"
	"ne-project/internal/services/product"
)

type Services struct {
	Category category.CategoryServiceItf
	Product product.ProductServiceItf
}

func NewServices(repositories *repositories.Repositories) *Services {
	return &Services{
		Category: category.NewCategoryService(
			repositories.Category,
		),
		Product: product.NewProductService(
			repositories.Product,
		),
	}
}