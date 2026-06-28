package entity

import (
	"encoding/json"
	"time"
)

type Review struct {
	ID              string          `db:"id" json:"id"`
	ProductID       string          `db:"product_id" json:"product_id"`
	BuyerID         string          `db:"buyer_id" json:"buyer_id"`
	OrderID         string          `db:"order_id" json:"order_id"`
	Rating          int             `db:"rating" json:"rating"`
	Comment         *string         `db:"comment" json:"comment"`
	Images          json.RawMessage `db:"images" json:"images"`
	SellerReply     *string         `db:"seller_reply" json:"seller_reply"`
	SellerRepliedAt *time.Time      `db:"seller_replied_at" json:"seller_replied_at"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
}

type ReviewWithBuyer struct {
	Review
	BuyerName   string  `db:"buyer_name" json:"buyer_name"`
	BuyerAvatar *string `db:"buyer_avatar" json:"buyer_avatar"`
}
