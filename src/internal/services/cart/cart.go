package cart

import (
	"context"

	"ne-project/src/internal/models/dto"
)

type CartServiceItf interface {
	AddToCart(ctx context.Context, buyerID string, req *dto.AddToCartRequest) (*dto.CartItemResponse, error)
	GetCart(ctx context.Context, buyerID string) (*dto.CartResponse, error)
	UpdateQuantity(ctx context.Context, buyerID string, cartID string, req *dto.UpdateCartRequest) error
	RemoveFromCart(ctx context.Context, buyerID string, cartID string) error
	ClearCart(ctx context.Context, buyerID string) error
}
