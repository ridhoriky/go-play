package wishlist

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"

	"github.com/jmoiron/sqlx"
)

type WishlistRepositoryItf interface {
	Add(ctx context.Context, item entity.Wishlist) error
	Remove(ctx context.Context, buyerID, productID string) error
	GetByBuyerID(ctx context.Context, buyerID string, params *dto.GetWishlistQuery) ([]entity.Wishlist, int, error)
	IsInWishlist(ctx context.Context, buyerID, productID string) (bool, error)
}

type wishlistRepository struct {
	db *sqlx.DB
}

func NewWishlistRepository(db *sqlx.DB) WishlistRepositoryItf {
	return &wishlistRepository{
		db: db,
	}
}
