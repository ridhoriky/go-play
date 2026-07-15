package wishlist

import (
	"context"
	"net/http"
	"sync"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/repositories/product"
	"ne-project/src/internal/repositories/wishlist"
)

type wishlistService struct {
	wishlistRepo wishlist.WishlistRepositoryItf
	productRepo  product.ProductRepositoryItf
}

func NewWishlistService(
	wishlistRepo wishlist.WishlistRepositoryItf,
	productRepo product.ProductRepositoryItf,
) WishlistServiceItf {
	return &wishlistService{
		wishlistRepo: wishlistRepo,
		productRepo:  productRepo,
	}
}

func (s *wishlistService) AddToWishlist(ctx context.Context, buyerID, productID string) error {
	// Validate product
	prodDetail, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	if prodDetail == nil {
		return dto.NewError(http.StatusNotFound, "product not found")
	}
	if !prodDetail.IsActive {
		return dto.NewError(http.StatusBadRequest, "product is not active")
	}

	// Check if already in wishlist
	exists, err := s.wishlistRepo.IsInWishlist(ctx, buyerID, productID)
	if err != nil {
		return err
	}
	if exists {
		return dto.NewError(http.StatusConflict, "product is already in wishlist")
	}

	return s.wishlistRepo.Add(ctx, entity.Wishlist{
		BuyerID:   buyerID,
		ProductID: productID,
	})
}

func (s *wishlistService) RemoveFromWishlist(ctx context.Context, buyerID, productID string) error {
	return s.wishlistRepo.Remove(ctx, buyerID, productID)
}

func (s *wishlistService) ToggleWishlist(ctx context.Context, buyerID, productID string) (bool, error) {
	// Validate product
	prodDetail, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return false, err
	}
	if prodDetail == nil {
		return false, dto.NewError(http.StatusNotFound, "product not found")
	}
	if !prodDetail.IsActive {
		return false, dto.NewError(http.StatusBadRequest, "product is not active")
	}

	exists, err := s.wishlistRepo.IsInWishlist(ctx, buyerID, productID)
	if err != nil {
		return false, err
	}

	if exists {
		err = s.wishlistRepo.Remove(ctx, buyerID, productID)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	err = s.wishlistRepo.Add(ctx, entity.Wishlist{
		BuyerID:   buyerID,
		ProductID: productID,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *wishlistService) GetWishlist(ctx context.Context, buyerID string, params *dto.GetWishlistQuery) (*dto.WishlistResponse, error) {
	wishlists, total, err := s.wishlistRepo.GetByBuyerID(ctx, buyerID, params)
	if err != nil {
		return nil, err
	}

	items := make([]dto.WishlistItemResponse, len(wishlists))
	var wg sync.WaitGroup

	for i, w := range wishlists {
		wg.Add(1)
		go func(i int, w entity.Wishlist) {
			defer wg.Done()
			prodDetail, err := s.productRepo.GetByID(ctx, w.ProductID)
			if err != nil || prodDetail == nil {
				return // skip if not found or error
			}

			items[i] = dto.WishlistItemResponse{
				ID: w.ID,
				Product: dto.ProductResponse{
					ID:           prodDetail.ID,
					StoreID:      prodDetail.StoreID,
					StoreName:    prodDetail.StoreName,
					StoreSlug:    prodDetail.StoreSlug,
					Name:         prodDetail.Name,
					Slug:         prodDetail.Slug,
					Description:  prodDetail.Description,
					Price:        prodDetail.Price,
					Stock:        prodDetail.Stock,
					CategoryID:   prodDetail.CategoryID,
					CategoryName: prodDetail.CategoryName,
					RatingAvg:    prodDetail.RatingAvg,
					TotalSold:    prodDetail.TotalSold,
					IsActive:     prodDetail.IsActive,
					IsInWishlist: true,
					PrimaryImage: prodDetail.PrimaryImage,
					CreatedAt:    prodDetail.CreatedAt,
					UpdatedAt:    prodDetail.UpdatedAt,
				},
				AddedAt: w.CreatedAt,
			}
		}(i, w)
	}
	wg.Wait()

	finalItems := make([]dto.WishlistItemResponse, 0, len(items))
	for i := range items {
		if items[i].ID != "" {
			finalItems = append(finalItems, items[i])
		}
	}

	return &dto.WishlistResponse{
		Items:      finalItems,
		TotalItems: total,
	}, nil
}
