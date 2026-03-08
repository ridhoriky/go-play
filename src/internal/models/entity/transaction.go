package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	ID          string              `db:"id" json:"id"`
	TotalAmount decimal.Decimal     `db:"total_amount" json:"total_amount"`
	CreatedAt   time.Time           `db:"created_at" json:"created_at"`
	Details     []TransactionDetail `json:"details"`
}
