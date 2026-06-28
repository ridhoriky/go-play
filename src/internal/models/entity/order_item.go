package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type OrderItem struct {
	ID           string          `db:"id" json:"id"`
	OrderID      string          `db:"order_id" json:"order_id"`
	ProductID    string          `db:"product_id" json:"product_id"`
	ProductName  string          `db:"product_name" json:"product_name"`
	ProductImage *string         `db:"product_image" json:"product_image"`
	Quantity     int             `db:"quantity" json:"quantity"`
	Price        decimal.Decimal `db:"price" json:"price"`
	Subtotal     decimal.Decimal `db:"subtotal" json:"subtotal"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
}
