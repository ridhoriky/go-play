package entity

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID              string          `db:"id" json:"id"`
	BuyerID         string          `db:"buyer_id" json:"buyer_id"`
	StoreID         string          `db:"store_id" json:"store_id"`
	OrderNumber     string          `db:"order_number" json:"order_number"`
	TotalAmount     decimal.Decimal `db:"total_amount" json:"total_amount"`
	Status          string          `db:"status" json:"status"`
	ShippingAddress json.RawMessage `db:"shipping_address" json:"shipping_address"`
	ShippingCost    decimal.Decimal `db:"shipping_cost" json:"shipping_cost"`
	PaymentMethod   *string         `db:"payment_method" json:"payment_method"`
	PaymentRef      *string         `db:"payment_ref" json:"payment_ref"`
	Notes           *string         `db:"notes" json:"notes"`
	CreatedAt       time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at" json:"updated_at"`
}
