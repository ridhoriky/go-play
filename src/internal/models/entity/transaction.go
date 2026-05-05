package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionStatus string

type Transaction struct {
	ID          string            `db:"id" json:"id"`
	TotalAmount decimal.Decimal   `db:"total_amount" json:"total_amount"`
	Status      TransactionStatus `db:"status" json:"status"`
	CreatedAt   time.Time         `db:"created_at" json:"created_at"`
}

type TransactionDetail struct {
	ID            string          `db:"id" json:"id"`
	TransactionID string          `db:"transaction_id" json:"transaction_id"`
	ProductID     string          `db:"product_id" json:"product_id"`
	ProductName   string          `db:"product_name" json:"product_name"`
	Quantity      int             `db:"quantity" json:"quantity"`
	Price         decimal.Decimal `db:"price" json:"price"`
	Subtotal      decimal.Decimal `db:"subtotal" json:"subtotal"`
}

type TransactionWithDetails struct {
	Transaction
	Items []TransactionDetail `json:"items"`
}

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusPaid      TransactionStatus = "paid"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)
