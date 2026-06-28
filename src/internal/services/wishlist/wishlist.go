package wishlist

import (
	"context"

	"ne-project/src/internal/models/dto"
)

type WishlistServiceItf interface {
	AddToWishlist(ctx context.Context, buyerID, productID string) error
	RemoveFromWishlist(ctx context.Context, buyerID, productID string) error
	GetWishlist(ctx context.Context, buyerID string, params *dto.GetWishlistQuery) (*dto.WishlistResponse, error)
	ToggleWishlist(ctx context.Context, buyerID, productID string) (bool, error)
}
