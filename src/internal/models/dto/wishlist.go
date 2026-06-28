package dto

import "time"

type AddToWishlistRequest struct {
	ProductID string `json:"product_id" binding:"required,uuid"`
}

type WishlistItemResponse struct {
	ID      string          `json:"id"`
	Product ProductResponse `json:"product"`
	AddedAt time.Time       `json:"added_at"`
}

type WishlistResponse struct {
	Items      []WishlistItemResponse `json:"items"`
	TotalItems int                    `json:"total_items"`
}

type GetWishlistQuery struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}
