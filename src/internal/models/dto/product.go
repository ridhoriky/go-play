package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

// ─── Request ────────────────────────────────────────────────────────────────

type CreateProductRequest struct {
	CategoryID  string  `json:"category_id" binding:"required,uuid"`
	Name        string  `json:"name"        binding:"required,min=1,max=255"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price"       binding:"required,min=0"`
	Stock       int     `json:"stock"       binding:"required,min=0"`
}

type UpdateProductRequest struct {
	CategoryID  string  `json:"category_id" binding:"required,uuid"`
	Name        string  `json:"name"        binding:"required,min=1,max=255"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price"       binding:"required,min=0"`
	Stock       int     `json:"stock"       binding:"required,min=0"`
}

type UpdateProductStockRequest struct {
	Stock int `json:"stock" binding:"required,min=0"`
}

type AddProductImageRequest struct {
	URL     string `json:"url" binding:"required,url"`
	AltText string `json:"alt_text"`
}

type CreateMultipleProducts struct {
	Data []*CreateProductRequest `json:"data" binding:"required,min=1,dive"`
}

type GetProductsQuery struct {
	Page     int             `form:"page"        binding:"omitempty,min=1"`
	Limit    int             `form:"limit"       binding:"omitempty,min=1,max=100"`
	Q        string          `form:"q"           binding:"omitempty,max=100"`
	Category string          `form:"category"    binding:"omitempty,uuid"`
	Store    string          `form:"store"       binding:"omitempty,uuid"`
	MinPrice decimal.Decimal `form:"min_price"   binding:"omitempty,min=0"`
	MaxPrice decimal.Decimal `form:"max_price"   binding:"omitempty,min=0"`
	Rating   float64         `form:"rating"      binding:"omitempty,min=0,max=5"`
	InStock  bool            `form:"in_stock"    binding:"omitempty"`
	LowStock bool            `form:"low_stock"   binding:"omitempty"`
	IsActive *bool           `form:"is_active"   binding:"omitempty"`
	Sort     string          `form:"sort"        binding:"omitempty,oneof=newest price_asc price_desc rating popular"`
}

// ─── Response ────────────────────────────────────────────────────────────────

type CategorySummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductImageResponse struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	AltText   string    `json:"alt_text"`
	SortOrder int       `json:"sort_order"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
}

type ProductResponse struct {
	ID              string          `json:"id"`
	StoreID         string          `json:"store_id"`
	StoreName       string          `json:"store_name,omitempty"`
	StoreSlug       string          `json:"store_slug,omitempty"`
	StoreIsVerified bool            `json:"store_is_verified"`
	Name            string          `json:"name"`
	Slug            string          `json:"slug"`
	Description     string          `json:"description"`
	Price           decimal.Decimal `json:"price"`
	Stock           int             `json:"stock"`
	CategoryID      string          `json:"category_id"`
	CategoryName    string          `json:"category_name"`
	RatingAvg       float64         `json:"rating_avg"`
	TotalSold       int             `json:"total_sold"`
	IsActive        bool            `json:"is_active"`
	IsInWishlist    bool            `json:"is_in_wishlist,omitempty"`
	PrimaryImage    string          `json:"primary_image,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type ProductStoreSummary struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Slug       string  `json:"slug"`
	IsVerified bool    `json:"is_verified"`
	LogoURL    string  `json:"logo_url"`
	RatingAvg  float64 `json:"rating_avg"`
}

type ProductDetailResponse struct {
	ID           string                 `json:"id"`
	Store        ProductStoreSummary    `json:"store"`
	Name         string                 `json:"name"`
	Slug         string                 `json:"slug"`
	Description  string                 `json:"description"`
	Price        decimal.Decimal        `json:"price"`
	Stock        int                    `json:"stock"`
	Category     CategorySummary        `json:"category"`
	RatingAvg    float64                `json:"rating_avg"`
	TotalReviews int                    `json:"total_reviews"`
	TotalSold    int                    `json:"total_sold"`
	IsActive     bool                   `json:"is_active"`
	IsInWishlist bool                   `json:"is_in_wishlist,omitempty"`
	Images       []ProductImageResponse `json:"images"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type ProductListResponse struct {
	Data []ProductResponse `json:"data"`
	Meta PaginationMeta    `json:"meta"`
}

type SellerProductDetailResponse struct {
	ProductDetailResponse
	RecentOrders []OrderSummary `json:"recent_orders"`
}
