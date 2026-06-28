package dto

import "time"

type CreateStoreRequest struct {
	StoreName   string `json:"store_name" binding:"required,min=3"`
	Description string `json:"description" binding:"omitempty"`
}

type UpdateStoreRequest struct {
	StoreName   string `json:"store_name" binding:"required,min=3"`
	Description string `json:"description" binding:"omitempty"`
	LogoURL     string `json:"logo_url" binding:"omitempty,url"`
	BannerURL   string `json:"banner_url" binding:"omitempty,url"`
}

type StoreResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	StoreName   string    `json:"store_name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	LogoURL     string    `json:"logo_url"`
	BannerURL   string    `json:"banner_url"`
	IsVerified  bool      `json:"is_verified"`
	RatingAvg   float64   `json:"rating_avg"`
	TotalSales  int       `json:"total_sales"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetStoresQuery struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search string `form:"search" binding:"omitempty,max=100"`
}

type StoreListResponse struct {
	Data []StoreResponse `json:"data"`
	Meta PaginationMeta  `json:"meta"`
}
