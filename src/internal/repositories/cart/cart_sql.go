package cart

import (
	"context"
	"database/sql"
	"errors"

	"ne-project/src/internal/models/entity"

	"github.com/jmoiron/sqlx"
)

type cartRepository struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) CartRepositoryItf {
	return &cartRepository{
		db: db,
	}
}

func (r *cartRepository) Add(ctx context.Context, item *entity.Cart) error {
	_, err := r.db.NamedExecContext(ctx, addCartQuery, item)
	return err
}

func (r *cartRepository) GetByBuyerID(ctx context.Context, buyerID string) ([]entity.Cart, error) {
	var carts []entity.Cart
	err := r.db.SelectContext(ctx, &carts, getCartByBuyerIDQuery, buyerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []entity.Cart{}, nil
		}
		return nil, err
	}
	return carts, nil
}

func (r *cartRepository) GetByID(ctx context.Context, id string) (*entity.Cart, error) {
	var cart entity.Cart
	err := r.db.GetContext(ctx, &cart, getCartByIDQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) UpdateQuantity(ctx context.Context, id string, quantity int) error {
	_, err := r.db.ExecContext(ctx, updateCartQuantityQuery, quantity, id)
	return err
}

func (r *cartRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, deleteCartQuery, id)
	return err
}

func (r *cartRepository) DeleteByBuyerID(ctx context.Context, buyerID string) error {
	_, err := r.db.ExecContext(ctx, deleteCartByBuyerIDQuery, buyerID)
	return err
}

func (r *cartRepository) GetByBuyerAndProduct(ctx context.Context, buyerID string, productID string) (*entity.Cart, error) {
	var cart entity.Cart
	err := r.db.GetContext(ctx, &cart, getCartByBuyerAndProductQuery, buyerID, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return &cart, nil
}
