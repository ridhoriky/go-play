package entity

import "time"

type Cart struct {
	ID        string    `db:"id" json:"id"`
	BuyerID   string    `db:"buyer_id" json:"buyer_id"`
	ProductID string    `db:"product_id" json:"product_id"`
	Quantity  int       `db:"quantity" json:"quantity"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
