package entity

import "github.com/shopspring/decimal"

type TransactionDetail struct {
	ID            string          `db:"id" json:"id"`
	TransactionID string          `db:"transaction_id" json:"transaction_id"`
	ProductID     string          `db:"product_id" json:"product_id"`
	ProductName   string          `db:"product_name" json:"product_name,omitempty"`
	Quantity      int             `db:"quantity" json:"quantity"`
	Price         decimal.Decimal `db:"price" json:"price"`
	Subtotal      decimal.Decimal `db:"subtotal" json:"subtotal"`
}
