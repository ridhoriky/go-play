package product

import (
	"context"
	"math"
	"net/http"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"

	"github.com/shopspring/decimal"
)

func (s *productService) GetAllProducts(ctx context.Context, req *dto.GetProductsQuery) (*dto.ProductListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	products, total, err := s.productRepository.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	res := make([]entity.ProductWithCategory, 0, len(products))

	for _, p := range products {
		resProduct := entity.ProductWithCategory{
			Product: entity.Product{
				ID:        p.ID,
				Name:      p.Name,
				Price:     p.Price,
				Stock:     p.Stock,
				CreatedAt: p.CreatedAt,
				UpdatedAt: p.UpdatedAt,
			},
			CategoryName: p.CategoryName,
		}

		res = append(res, resProduct)
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	resp := &dto.ProductListResponse{
		Data: res,
		Meta: dto.PaginationMeta{
			Total:      total,
			Page:       req.Page,
			Limit:      req.Limit,
			TotalPages: totalPages,
		},
	}

	return resp, nil
}

func (s *productService) CreateProduct(ctx context.Context, product *dto.CreateProductRequest) (*entity.Product, error) {

	if product.Price < 0 {
		return nil, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrProductPriceNegative,
		}
	}

	if product.Stock < 0 {
		return nil, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrProductStockNegative,
		}
	}

	if product.Name == "" {
		return nil, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrProductNameRequied,
		}
	}

	p := entity.Product{
		Name:       product.Name,
		Price:      decimal.NewFromFloat(product.Price),
		Stock:      product.Stock,
		CategoryID: product.CategoryID,
	}
	if err := s.productRepository.Create(ctx, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *productService) GetProductByID(ctx context.Context, id string) (*dto.ProductResponse, error) {
	product, CategoryName, err := s.productRepository.GetByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return &dto.ProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Price:        product.Price,
		Stock:        product.Stock,
		CategoryID:   product.CategoryID,
		CategoryName: CategoryName,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}, nil

}

func (s *productService) UpdateProduct(ctx context.Context, id string, req *dto.UpdateProductRequest) (*entity.Product, error) {
	existingProduct, _, err := s.productRepository.GetByID(ctx, id)

	if err != nil {
		return nil, &dto.Error{
			Code:    http.StatusNotFound,
			Message: preference.ErrInvalidProductID,
		}
	}

	existingProduct.Name = req.Name
	existingProduct.Price = decimal.NewFromFloat(req.Price)
	existingProduct.Stock = req.Stock
	existingProduct.CategoryID = req.CategoryID

	if err := s.productRepository.Update(ctx, id, existingProduct); err != nil {
		return nil, err
	}

	return existingProduct, err
}

func (s *productService) DeleteProduct(ctx context.Context, id string) error {
	return s.productRepository.Delete(ctx, id)
}

func (s *productService) CreateMultipleProducts(ctx context.Context, products []entity.Product) ([]entity.Product, error) {

	return s.productRepository.CreateMultiple(ctx, products)
}
