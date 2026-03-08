package product

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
)

func (s *productService) GetAllProducts(ctx context.Context, req dto.ProductFilterRequest) ([]dto.ProductResponse, error) {
	if req.Limit > 100 {
		req.Limit = 100
	}

	if req.Page < 1 {
		req.Page = 1
	}
	return s.productRepository.GetAll(ctx, req)
}

func (s *productService) CreateProduct(ctx context.Context, product *dto.ProductDTO) error {
	return s.productRepository.Create(ctx, product.ToModelPtr())
}

func (s *productService) GetProductByID(ctx context.Context, id string) (*dto.ProductResponse, error) {
	return s.productRepository.GetByID(ctx, id)
}

func (s *productService) UpdateProduct(ctx context.Context, product *dto.ProductDTO) error {
	return s.productRepository.Update(ctx, product.ToModelPtr())
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	return s.productRepository.Delete(ctx, id)
}

func (s *productService) CreateMultipleProducts(ctx context.Context, products []dto.ProductDTO) ([]dto.ProductDTO, error) {
	models := make([]entity.Product, len(products))
	for i, p := range products {
		models[i] = *p.ToModelPtr()
	}

	return s.productRepository.CreateMultiple(ctx, models)
}
