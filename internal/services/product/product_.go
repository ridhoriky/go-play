package product

import (
	"context"
	"ne-project/internal/dto"
	"ne-project/internal/models"
)

func (s *productService) GetAllProducts(ctx context.Context) ([]dto.ProductResponse, error) {
	return s.productRepository.GetAll(ctx)
}

func (s *productService) CreateProduct(ctx context.Context, product *dto.ProductDTO) error {
	return s.productRepository.Create(ctx, product.ToModelPtr())
}

func (s *productService) GetProductByID(ctx context.Context, id int) (*dto.ProductResponse, error) {
	return s.productRepository.GetByID(ctx, id)
}

func (s *productService) UpdateProduct(ctx context.Context, product *dto.ProductDTO) error {
	return s.productRepository.Update(ctx, product.ToModelPtr())
}

func (s *productService) DeleteProduct(ctx context.Context, id int) error {
	return s.productRepository.Delete(ctx, id)
}

func (s *productService) CreateMultipleProducts(ctx context.Context, products []dto.ProductDTO) ([]dto.ProductResponse, error) {
	models := make([]models.Product, len(products))
	for i, p := range products {
		models[i] = *p.ToModelPtr()
	}
	
	return s.productRepository.CreateMultiple(ctx, models)
}