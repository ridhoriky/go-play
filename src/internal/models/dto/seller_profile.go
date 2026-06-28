package dto

import (
	"github.com/shopspring/decimal"
)

type UpdateSellerProfileRequest struct {
	Name        string `json:"name" binding:"omitempty"`
	Phone       string `json:"phone" binding:"omitempty"`
	AvatarURL   string `json:"avatar_url" binding:"omitempty,url"`
	StoreName   string `json:"store_name" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	LogoURL     string `json:"logo_url" binding:"omitempty,url"`
	BannerURL   string `json:"banner_url" binding:"omitempty,url"`
}

type SellerProfileResponse struct {
	User  UserResponse  `json:"user"`
	Store StoreResponse `json:"store"`
}

type StoreStats struct {
	TotalProducts  int             `json:"total_products"`
	ActiveProducts int             `json:"active_products"`
	TotalOrders    int             `json:"total_orders"`
	TotalRevenue   decimal.Decimal `json:"total_revenue"`
	AverageRating  float64         `json:"average_rating"`
	TotalReviews   int             `json:"total_reviews"`
}
