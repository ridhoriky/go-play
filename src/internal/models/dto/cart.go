package dto

import "github.com/shopspring/decimal"

type AddToCartRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid4"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type UpdateCartRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

type CartItemResponse struct {
	ID       string          `json:"id"`
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity"`
	Subtotal decimal.Decimal `json:"subtotal"`
	Store    StoreResponse   `json:"store"`
}

type CartResponse struct {
	Items       []CartItemResponse `json:"items"`
	TotalAmount decimal.Decimal    `json:"total_amount"`
	TotalItems  int                `json:"total_items"`
}
