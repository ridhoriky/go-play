package dto

import "time"

// ─── Request ────────────────────────────────────────────────────────────────

type CreateCategoryRequest struct {
	Name        string  `json:"name"        binding:"required,min=1,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	ParentID    *string `json:"parent_id"   binding:"omitempty,uuid"`
	ImageURL    *string `json:"image_url"   binding:"omitempty,url,max=512"`
	SortOrder   int     `json:"sort_order"  binding:"omitempty,min=0"`
}

type UpdateCategoryRequest struct {
	Name        string  `json:"name"        binding:"required,min=1,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	ParentID    *string `json:"parent_id"   binding:"omitempty,uuid"`
	ImageURL    *string `json:"image_url"   binding:"omitempty,url,max=512"`
	SortOrder   int     `json:"sort_order"  binding:"omitempty,min=0"`
}

// ─── Response ────────────────────────────────────────────────────────────────

type CategoryResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	ParentID    *string    `json:"parent_id,omitempty"`
	ImageURL    *string    `json:"image_url,omitempty"`
	SortOrder   int        `json:"sort_order"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type CategoryDetailResponse struct {
	CategoryResponse
	ProductsCount int                `json:"products_count"`
	Children      []CategoryResponse `json:"children"`
	Breadcrumb    []CategoryResponse `json:"breadcrumb"`
}

type CategoryTreeNode struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Description  *string            `json:"description"`
	ImageURL     *string            `json:"image_url"`
	SortOrder    int                `json:"sort_order"`
	ProductCount int                `json:"product_count"`
	Children     []CategoryTreeNode `json:"children"`
}

type DeleteCategoryResponse struct {
	Message   string    `json:"message"`
	DeletedAt time.Time `json:"deleted_at"`
}
