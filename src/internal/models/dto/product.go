package dto

import (
	"ne-project/src/internal/models/entity"
	"time"

	"github.com/shopspring/decimal"
)

// ─── Request ────────────────────────────────────────────────────────────────

type CreateProductRequest struct {
	CategoryID string  `json:"category_id" binding:"required,uuid"`
	Name       string  `json:"name"        binding:"required,min=1,max=255"`
	Price      float64 `json:"price"       binding:"required,min=0"`
	Stock      int     `json:"stock"       binding:"required,min=0"`
}

type UpdateProductRequest struct {
	CategoryID string  `json:"category_id" binding:"required,uuid"`
	Name       string  `json:"name"        binding:"required,min=1,max=255"`
	Price      float64 `json:"price"       binding:"required,min=0"`
	Stock      int     `json:"stock"       binding:"required,min=0"`
}

type UpdateProductStockRequest struct {
	Stock int `json:"stock" binding:"required,min=0"`
}

type GetProductsQuery struct {
	Page       int             `form:"page"        binding:"omitempty,min=1"`
	Limit      int             `form:"limit"       binding:"omitempty,min=1,max=100"`
	Search     string          `form:"search"      binding:"omitempty,max=100"`
	CategoryID string          `form:"category_id" binding:"omitempty,uuid"`
	MinPrice   decimal.Decimal `form:"min_price"   binding:"omitempty,min=0"`
	MaxPrice   decimal.Decimal `form:"max_price"   binding:"omitempty,min=0"`
	InStock    bool            `form:"in_stock"    binding:"omitempty"`
	SortBy     string          `form:"sort_by"     binding:"omitempty,oneof=name price stock created_at"`
	SortDir    string          `form:"sort_dir"    binding:"omitempty,oneof=asc desc"`
}

// ─── Response ────────────────────────────────────────────────────────────────

type CategorySummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductResponse struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Price        decimal.Decimal `json:"price"`
	Stock        int             `json:"stock"`
	CategoryID   string          `json:"category_id"`
	CategoryName string          `json:"category_name"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

type ProductDetailResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Price     decimal.Decimal `json:"price"`
	Stock     int             `json:"stock"`
	Category  CategorySummary `json:"category"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ProductListResponse struct {
	Data []entity.ProductWithCategory `json:"data"`
	Meta PaginationMeta               `json:"meta"`
}
