package cart

import (
	"context"

	"ne-project/src/internal/models/entity"
)

type CartRepositoryItf interface {
	Add(ctx context.Context, item *entity.Cart) error
	GetByBuyerID(ctx context.Context, buyerID string) ([]entity.Cart, error)
	GetByID(ctx context.Context, id string) (*entity.Cart, error)
	UpdateQuantity(ctx context.Context, id string, quantity int) error
	Delete(ctx context.Context, id string) error
	DeleteByBuyerID(ctx context.Context, buyerID string) error
	GetByBuyerAndProduct(ctx context.Context, buyerID string, productID string) (*entity.Cart, error)
}
