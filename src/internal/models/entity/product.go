package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ID         string          `db:"id" json:"id"`
	CategoryID string          `db:"category_id" json:"category_id"`
	Name       string          `db:"name" json:"name"`
	Price      decimal.Decimal `db:"price" json:"price"`
	Stock      int             `db:"stock" json:"stock"`
	CreatedAt  time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time       `db:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time      `db:"deleted_at" json:"deleted_at,omitempty"`
}

type ProductWithCategory struct {
	Product
	CategoryName string `db:"category_name" json:"category_name"`
}
