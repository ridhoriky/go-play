package dto

import "time"

// ─── Request ────────────────────────────────────────────────────────────────

type CreateCategoryRequest struct {
	Name        string  `json:"name"        binding:"required,min=1,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
}

type UpdateCategoryRequest struct {
	Name        string  `json:"name"        binding:"required,min=1,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
}

type GetCategoriesQuery struct {
	Page           int    `form:"page"            binding:"omitempty,min=1"`
	Limit          int    `form:"limit"           binding:"omitempty,min=1,max=100"`
	Search         string `form:"search"          binding:"omitempty,max=100"`
	IncludeDeleted bool   `form:"include_deleted" binding:"omitempty"`
	SortBy         string `form:"sort_by"     binding:"omitempty,oneof=name created_at"`
	SortDir        string `form:"sort_dir"    binding:"omitempty,oneof=asc desc"`
}

// ─── Response ────────────────────────────────────────────────────────────────

type CategoryResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

type CategoryDetailResponse struct {
	CategoryResponse
	ProductsCount int `json:"products_count"`
}

type CategoryListResponse struct {
	Data []CategoryResponse `json:"data"`
	Meta PaginationMeta     `json:"meta"`
}

type DeleteCategoryResponse struct {
	Message   string    `json:"message"`
	DeletedAt time.Time `json:"deleted_at"`
}
