package product

import (
	"context"
	"ne-project/internal/dto"
	"ne-project/internal/repositories/product"
)

type ProductServiceItf interface {
	GetAllProducts(ctx context.Context, req dto.ProductFilterRequest) ([]dto.ProductResponse, error)
	GetProductByID(ctx context.Context, id int) (*dto.ProductResponse, error)
	CreateProduct(ctx context.Context, product *dto.ProductDTO) error
	UpdateProduct(ctx context.Context, product *dto.ProductDTO) error
	DeleteProduct(ctx context.Context, id int) error
	CreateMultipleProducts(ctx context.Context, products []dto.ProductDTO) ([]dto.ProductDTO, error)
}

	
type productService struct {
	productRepository product.ProductRepositoryItf
}

func NewProductService(productRepository product.ProductRepositoryItf) ProductServiceItf {
	return &productService{
		productRepository: productRepository,
	}
}
