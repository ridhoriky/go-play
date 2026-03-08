package entity

import "github.com/shopspring/decimal"

type Product struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Price      decimal.Decimal `json:"price"`
	Stock      int             `json:"stock"`
	CategoryID string          `json:"category_id"`
}
