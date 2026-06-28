package dto

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type CheckoutRequest struct {
	CartIDs         []string        `json:"cart_ids" validate:"required,min=1"`
	ShippingAddress json.RawMessage `json:"shipping_address" validate:"required"`
	Notes           string          `json:"notes"`
}

type OrderResponse struct {
	ID          string             `json:"id"`
	OrderNumber string             `json:"order_number"`
	Store       StoreSummary       `json:"store"`
	Items       []OrderItemSummary `json:"items"`
	TotalAmount decimal.Decimal    `json:"total_amount"`
	Status      string             `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
}

type StoreSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type OrderItemSummary struct {
	ID           string          `json:"id"`
	ProductID    string          `json:"product_id"`
	ProductName  string          `json:"product_name"`
	ProductImage *string         `json:"product_image"`
	Quantity     int             `json:"quantity"`
	Price        decimal.Decimal `json:"price"`
	Subtotal     decimal.Decimal `json:"subtotal"`
}

type OrderDetailResponse struct {
	OrderResponse
	ShippingAddress json.RawMessage `json:"shipping_address"`
	ShippingCost    decimal.Decimal `json:"shipping_cost"`
	PaymentMethod   *string         `json:"payment_method"`
	Notes           *string         `json:"notes"`
}

type OrderListResponse struct {
	Data []OrderResponse `json:"data"`
	Meta PaginationMeta  `json:"meta"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=paid processing shipped delivered canceled refunded"`
}

type GetOrdersQuery struct {
	Page    int    `form:"page,default=1"`
	Limit   int    `form:"limit,default=10"`
	Status  string `form:"status"`
	SortBy  string `form:"sort_by,default=created_at"`
	SortDir string `form:"sort_dir,default=DESC"`
}
