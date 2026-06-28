package wishlist

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
)

func (r *wishlistRepository) Add(ctx context.Context, item entity.Wishlist) error {
	stmt, err := r.db.PrepareNamedContext(ctx, queryAddWishlist)
	if err != nil {
		return err
	}
	defer func() {
		_ = stmt.Close()
	}()

	err = stmt.QueryRowxContext(ctx, item).StructScan(&item)
	return err
}

func (r *wishlistRepository) Remove(ctx context.Context, buyerID, productID string) error {
	_, err := r.db.ExecContext(ctx, queryRemoveWishlist, buyerID, productID)
	return err
}

func (r *wishlistRepository) GetByBuyerID(ctx context.Context, buyerID string, params *dto.GetWishlistQuery) ([]entity.Wishlist, int, error) {
	page := 1
	limit := 10
	if params != nil {
		if params.Page > 0 {
			page = params.Page
		}
		if params.Limit > 0 {
			limit = params.Limit
		}
	}
	offset := (page - 1) * limit

	var items []entity.Wishlist
	err := r.db.SelectContext(ctx, &items, queryGetWishlistsByBuyerID, buyerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var count int
	err = r.db.GetContext(ctx, &count, queryCountWishlistsByBuyerID, buyerID)
	if err != nil {
		return nil, 0, err
	}

	return items, count, nil
}

func (r *wishlistRepository) IsInWishlist(ctx context.Context, buyerID, productID string) (bool, error) {
	var exists bool
	err := r.db.GetContext(ctx, &exists, queryCheckInWishlist, buyerID, productID)
	return exists, err
}
