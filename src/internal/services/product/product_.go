package product

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/utils"
	"ne-project/src/internal/utils/validation"

	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

func (s *productService) GetAllProducts(ctx context.Context, userID string, req *dto.GetProductsQuery) (*dto.ProductListResponse, error) {
	req.Page = validation.ValidatePage(req.Page)
	req.Limit = validation.ValidatePageSize(req.Limit)

	products, total, err := s.productRepository.GetAll(ctx, req)

	if err != nil {
		return nil, err
	}

	res := make([]dto.ProductResponse, 0, len(products))

	for i := range products {
		p := &products[i]

		isInWishlist := false
		if userID != "" {
			inWishlist, err := s.wishlistRepository.IsInWishlist(ctx, userID, p.ID)
			if err == nil {
				isInWishlist = inWishlist
			}
		}

		resProduct := dto.ProductResponse{
			ID:              p.ID,
			StoreID:         p.StoreID,
			StoreName:       p.StoreName,
			StoreSlug:       p.StoreSlug,
			StoreIsVerified: p.StoreIsVerified,
			PrimaryImage:    p.PrimaryImage,
			Name:            p.Name,
			Slug:            p.Slug,
			Description:     p.Description,
			Price:           p.Price,
			Stock:           p.Stock,
			CategoryID:      p.CategoryID,
			CategoryName:    p.CategoryName,
			RatingAvg:       p.RatingAvg,
			TotalSold:       p.TotalSold,
			IsActive:        p.IsActive,
			IsInWishlist:    isInWishlist,
			CreatedAt:       p.CreatedAt,
			UpdatedAt:       p.UpdatedAt,
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

func (s *productService) CreateProduct(ctx context.Context, storeID string, product *dto.CreateProductRequest) (*entity.Product, error) {

	if product.Price < 0 {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrProductPriceNegative)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrProductPriceNegative)
	}

	if product.Stock < 0 {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrProductStockNegative)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrProductStockNegative)
	}

	if product.Name == "" {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrProductNameRequired)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrProductNameRequired)
	}

	baseSlug := utils.GenerateSlug(product.Name)
	var finalSlug string
	for i := range 100 {
		candidate := baseSlug
		if i > 0 {
			candidate = fmt.Sprintf("%s-%d", baseSlug, i)
		}
		_, err := s.productRepository.GetBySlug(ctx, candidate)
		if err != nil {
			// Assuming error means not found
			finalSlug = candidate
			break
		}
	}

	if finalSlug == "" {
		return nil, errors.New("failed to generate unique slug")
	}

	p := entity.Product{
		StoreID:     storeID,
		Name:        product.Name,
		Slug:        finalSlug,
		Description: product.Description,
		Price:       decimal.NewFromFloat(product.Price),
		Stock:       product.Stock,
		CategoryID:  product.CategoryID,
		IsActive:    true,
	}

	if err := s.productRepository.Create(ctx, &p); err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *productService) GetProductByID(ctx context.Context, id string, userID string) (*dto.ProductDetailResponse, error) {
	product, err := s.productRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	images, err := s.productImageRepository.GetByProductID(ctx, id)
	if err != nil {
		// Log error but continue to return product details even if images fail
		zerolog.Ctx(ctx).Error().Err(err).Str("productID", id).Msg("failed to get product images")
	}

	var imageResponses []dto.ProductImageResponse
	for _, img := range images {
		imageResponses = append(imageResponses, dto.ProductImageResponse{
			ID:        img.ID,
			URL:       img.URL,
			AltText:   img.AltText,
			SortOrder: img.SortOrder,
			IsPrimary: img.IsPrimary,
			CreatedAt: img.CreatedAt,
		})
	}

	isInWishlist := false
	if userID != "" {
		inWishlist, err := s.wishlistRepository.IsInWishlist(ctx, userID, product.ID)
		if err == nil {
			isInWishlist = inWishlist
		}
	}

	return &dto.ProductDetailResponse{
		ID: product.ID,
		Store: dto.ProductStoreSummary{
			ID:         product.StoreID,
			Name:       product.StoreName,
			Slug:       product.StoreSlug,
			IsVerified: product.StoreIsVerified,
			LogoURL:    product.StoreLogoURL,
			RatingAvg:  product.StoreRatingAvg,
		},
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category: dto.CategorySummary{
			ID:   product.CategoryID,
			Name: product.CategoryName,
		},
		RatingAvg:    product.RatingAvg,
		TotalReviews: product.TotalReviews,
		TotalSold:    product.TotalSold,
		IsActive:     product.IsActive,
		IsInWishlist: isInWishlist,
		Images:       imageResponses,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}, nil

}

func (s *productService) GetProductBySlug(ctx context.Context, slug string, userID string) (*dto.ProductDetailResponse, error) {
	product, err := s.productRepository.GetDetailBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	images, err := s.productImageRepository.GetByProductID(ctx, product.ID)
	if err != nil {
		// Log error but continue to return product details even if images fail
		zerolog.Ctx(ctx).Error().Err(err).Str("productID", product.ID).Msg("failed to get product images")
	}

	var imageResponses []dto.ProductImageResponse
	for _, img := range images {
		imageResponses = append(imageResponses, dto.ProductImageResponse{
			ID:        img.ID,
			URL:       img.URL,
			AltText:   img.AltText,
			SortOrder: img.SortOrder,
			IsPrimary: img.IsPrimary,
			CreatedAt: img.CreatedAt,
		})
	}

	isInWishlist := false
	if userID != "" {
		inWishlist, err := s.wishlistRepository.IsInWishlist(ctx, userID, product.ID)
		if err == nil {
			isInWishlist = inWishlist
		}
	}

	return &dto.ProductDetailResponse{
		ID: product.ID,
		Store: dto.ProductStoreSummary{
			ID:         product.StoreID,
			Name:       product.StoreName,
			Slug:       product.StoreSlug,
			IsVerified: product.StoreIsVerified,
			LogoURL:    product.StoreLogoURL,
			RatingAvg:  product.StoreRatingAvg,
		},
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category: dto.CategorySummary{
			ID:   product.CategoryID,
			Name: product.CategoryName,
		},
		RatingAvg:    product.RatingAvg,
		TotalReviews: product.TotalReviews,
		TotalSold:    product.TotalSold,
		IsActive:     product.IsActive,
		IsInWishlist: isInWishlist,
		Images:       imageResponses,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id string, storeID string, req *dto.UpdateProductRequest) (*entity.Product, error) {
	existingProductDetail, err := s.productRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	existingProduct := &existingProductDetail.Product

	if existingProduct.StoreID != storeID {
		return nil, dto.NewError(http.StatusForbidden, "You do not have permission to update this product")
	}

	if req.Price < 0 {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrProductPriceNegative)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrProductPriceNegative)
	}

	if req.Stock < 0 {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrProductStockNegative)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrProductStockNegative)
	}

	if req.Name == "" {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrProductNameRequired)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrProductNameRequired)
	}

	existingProduct.Name = req.Name
	existingProduct.Description = req.Description
	existingProduct.Price = decimal.NewFromFloat(req.Price)
	existingProduct.Stock = req.Stock
	existingProduct.CategoryID = req.CategoryID
	// Slug is not typically updated, but could be. Assuming it stays the same.

	if err = s.productRepository.Update(ctx, id, existingProduct); err != nil {
		return nil, err
	}

	return existingProduct, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id string, storeID string) error {
	existingProductDetail, err := s.productRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	existingProduct := &existingProductDetail.Product

	if existingProduct.StoreID != storeID {
		return dto.NewError(http.StatusForbidden, "You do not have permission to delete this product")
	}

	hasActive, err := s.orderRepository.HasActiveOrdersForProduct(ctx, id)
	if err != nil {
		return err
	}
	if hasActive {
		return dto.NewError(http.StatusBadRequest, "cannot delete product with active orders")
	}

	return s.productRepository.Delete(ctx, id)
}

func (s *productService) CreateMultipleProducts(ctx context.Context, storeID string, products []entity.Product) ([]entity.Product, error) {
	for i := range products {
		products[i].StoreID = storeID
		products[i].IsActive = true

		baseSlug := utils.GenerateSlug(products[i].Name)
		var finalSlug string
		for j := range 100 {
			candidate := baseSlug
			if j > 0 {
				candidate = fmt.Sprintf("%s-%d", baseSlug, j)
			}
			_, err := s.productRepository.GetBySlug(ctx, candidate)
			if err != nil {
				finalSlug = candidate
				break
			}
		}

		if finalSlug == "" {
			return nil, fmt.Errorf("failed to generate unique slug for product %s", products[i].Name)
		}

		products[i].Slug = finalSlug
	}
	return s.productRepository.CreateMultiple(ctx, products)
}

func (s *productService) GetSellerProductDetail(ctx context.Context, id string, storeID string) (*dto.SellerProductDetailResponse, error) {
	productDetail, err := s.GetProductByID(ctx, id, "")
	if err != nil {
		return nil, err
	}

	if productDetail.Store.ID != storeID {
		return nil, dto.NewError(http.StatusForbidden, "You do not have permission to access this product")
	}

	recentOrders, err := s.orderRepository.GetRecentOrdersByProductID(ctx, id, 5)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to get recent orders")
		recentOrders = []dto.OrderSummary{}
	}

	if recentOrders == nil {
		recentOrders = []dto.OrderSummary{}
	}

	return &dto.SellerProductDetailResponse{
		ProductDetailResponse: *productDetail,
		RecentOrders:          recentOrders,
	}, nil
}

func (s *productService) AddProductImage(ctx context.Context, productID string, storeID string, req *dto.AddProductImageRequest) (*entity.ProductImage, error) {
	productDetail, err := s.productRepository.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if productDetail.StoreID != storeID {
		return nil, dto.NewError(http.StatusForbidden, "You do not have permission to update this product")
	}

	count, err := s.productImageRepository.CountImages(ctx, productID)
	if err != nil {
		return nil, err
	}
	if count >= 10 {
		return nil, dto.NewError(http.StatusBadRequest, "maximum 10 images allowed per product")
	}

	isPrimary := count == 0

	newImage := &entity.ProductImage{
		ProductID: productID,
		URL:       req.URL,
		AltText:   req.AltText,
		SortOrder: count,
		IsPrimary: isPrimary,
	}

	if err := s.productImageRepository.AddImage(ctx, newImage); err != nil {
		return nil, err
	}

	return newImage, nil
}

func (s *productService) DeleteProductImage(ctx context.Context, productID string, imageID string, storeID string) error {
	productDetail, err := s.productRepository.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	if productDetail.StoreID != storeID {
		return dto.NewError(http.StatusForbidden, "You do not have permission to update this product")
	}

	images, err := s.productImageRepository.GetByProductID(ctx, productID)
	if err != nil {
		return err
	}

	var targetImage *entity.ProductImage
	for _, img := range images {
		if img.ID == imageID {
			var captured = img
			targetImage = &captured
			break
		}
	}

	if targetImage == nil {
		return dto.NewError(http.StatusNotFound, "image not found")
	}

	if targetImage.IsPrimary && len(images) == 1 {
		return dto.NewError(http.StatusBadRequest, "can't delete primary image if it's the only one")
	}

	err = s.productImageRepository.DeleteImage(ctx, imageID)
	if err != nil {
		return err
	}

	if targetImage.IsPrimary && len(images) > 1 {
		for _, img := range images {
			if img.ID != imageID {
				_ = s.productImageRepository.SetPrimary(ctx, productID, img.ID)
				break
			}
		}
	}

	return nil
}

func (s *productService) SetPrimaryImage(ctx context.Context, productID string, imageID string, storeID string) error {
	productDetail, err := s.productRepository.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	if productDetail.StoreID != storeID {
		return dto.NewError(http.StatusForbidden, "You do not have permission to update this product")
	}

	images, err := s.productImageRepository.GetByProductID(ctx, productID)
	if err != nil {
		return err
	}

	imageExists := false
	for _, img := range images {
		if img.ID == imageID {
			imageExists = true
			break
		}
	}

	if !imageExists {
		return dto.NewError(http.StatusNotFound, "image not found")
	}

	if err := s.productImageRepository.UnsetPrimary(ctx, productID); err != nil {
		return err
	}

	return s.productImageRepository.SetPrimary(ctx, productID, imageID)
}
