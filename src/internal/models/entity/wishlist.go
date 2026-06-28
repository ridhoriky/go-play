package entity

import "time"

type Wishlist struct {
	ID        string    `db:"id"`
	BuyerID   string    `db:"buyer_id"`
	ProductID string    `db:"product_id"`
	CreatedAt time.Time `db:"created_at"`
}
