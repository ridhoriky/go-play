package dto

import "time"

// ─── Requests ─────────────────────────────────────────────────────

type CreateUserRequest struct {
	Name     string `json:"name"        binding:"required,min=1,max=100"`
	Email    string `json:"email"       binding:"required,email"`
	Role     string `json:"role"        binding:"required,oneof=admin user"`
	Password string `json:"password"    binding:"required,min=6"`
}

type UpdateUserRequest struct {
	Name      *string `json:"name"        binding:"omitempty,min=1,max=100"`
	Email     *string `json:"email"       binding:"omitempty,email"`
	Role      *string `json:"role"        binding:"omitempty,oneof=admin user"`
	Phone     *string `json:"phone"       binding:"omitempty"`
	AvatarURL *string `json:"avatar_url"  binding:"omitempty,url"`
}

type GetUsersQuery struct {
	Page     int    `form:"page"      binding:"omitempty,min=1"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search   string `form:"search"    binding:"omitempty,min=1,max=100"`
	IsActive *bool  `form:"is_active" binding:"omitempty"`
	SortBy   string `form:"sort_by"   binding:"omitempty,oneof=name email role created_at"`
	SortDir  string `form:"sort_dir"  binding:"omitempty,oneof=asc desc"`
}

// ─── Responses ────────────────────────────────────────────────────

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserListResponse struct {
	Data []UserResponse `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type DeleteUserResponse struct {
	Message   string    `json:"message"`
	DeletedAt time.Time `json:"deleted_at"`
}
