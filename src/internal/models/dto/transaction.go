package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// ─── Request ────────────────────────────────────────────────────────────────

type CheckoutItem struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
	Quantity  int    `json:"quantity"   binding:"required,min=1"`
}

type CreateTransactionRequest struct {
	Items []CheckoutItem `json:"items" binding:"required,min=1,dive"`
}

type UpdateTransactionStatusRequest struct {
	// oneof memastikan hanya nilai "paid" atau "cancelled" yang diterima
	Status string `json:"status" binding:"required,oneof=paid cancelled"`
}

type GetTransactionsQuery struct {
	Page     int    `form:"page"      binding:"omitempty,min=1"`
	Limit    int    `form:"limit"     binding:"omitempty,min=1,max=100"`
	Status   string `form:"status"    binding:"omitempty,oneof=pending paid cancelled"`
	DateFrom string `form:"date_from" binding:"omitempty"`
	DateTo   string `form:"date_to"   binding:"omitempty"`
	SortBy   string `form:"sort_by"   binding:"omitempty,oneof=created_at total_amount"`
	SortDir  string `form:"sort_dir"  binding:"omitempty,oneof=asc desc"`
}

// ─── Response ────────────────────────────────────────────────────────────────

type TransactionItemResponse struct {
	ID          string  `json:"id"`
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Subtotal    float64 `json:"subtotal"`
}

type TransactionResponse struct {
	ID          string    `json:"id"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	ItemsCount  int       `json:"items_count"`
	CreatedAt   time.Time `json:"created_at"`
}

type TransactionDetailResponse struct {
	ID          string                    `json:"id" example:"trx-123"`
	TotalAmount float64                   `json:"total_amount" example:"15000"`
	Status      string                    `json:"status" example:"paid"`
	CreatedAt   time.Time                 `json:"created_at" example:"2025-01-01T10:00:00Z"`
	Items       []TransactionItemResponse `json:"items"`
}

type TransactionListResponse struct {
	Data []TransactionResponse `json:"data"`
	Meta PaginationMeta        `json:"meta"`
}

type ProductSnapshot struct {
	ID    string
	Name  string
	Price decimal.Decimal
	Stock int
}
