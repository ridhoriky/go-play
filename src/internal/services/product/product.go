package product

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/repositories/product"
)

type ProductServiceItf interface {
	GetAllProducts(ctx context.Context, req *dto.GetProductsQuery) (*dto.ProductListResponse, error)
	GetProductByID(ctx context.Context, id string) (*dto.ProductResponse, error)
	CreateProduct(ctx context.Context, product *dto.CreateProductRequest) (*entity.Product, error)
	UpdateProduct(ctx context.Context, id string, product *dto.UpdateProductRequest) (*entity.Product, error)
	DeleteProduct(ctx context.Context, id string) error
	CreateMultipleProducts(ctx context.Context, products []entity.Product) ([]entity.Product, error)
}

type productService struct {
	productRepository product.ProductRepositoryItf
}

func NewProductService(productRepository product.ProductRepositoryItf) ProductServiceItf {
	return &productService{
		productRepository: productRepository,
	}
}
