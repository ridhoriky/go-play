package dto

import "time"

type AdminUserListParams struct {
	Search    string
	Role      string // admin, seller, buyer
	IsActive  *bool
	Page      int
	Limit     int
	SortBy    string // name, email, created_at
	SortOrder string // asc, desc
}

type AdminUserResponse struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	IsActive  bool       `json:"is_active"`
	HasStore  bool       `json:"has_store"`
	StoreName *string    `json:"store_name"`
	CreatedAt time.Time  `json:"created_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
}

type AdminUpdateUserRequest struct {
	Role     *string `json:"role" binding:"omitempty,oneof=buyer seller admin"`
	IsActive *bool   `json:"is_active" binding:"omitempty"`
}

type AdminSellerListParams struct {
	Search     string
	IsVerified *bool
	Page       int
	Limit      int
	SortBy     string
	SortOrder  string
}

type AdminSellerResponse struct {
	StoreID       string    `json:"store_id"`
	StoreName     string    `json:"store_name"`
	Slug          string    `json:"slug"`
	OwnerName     string    `json:"owner_name"`
	OwnerEmail    string    `json:"owner_email"`
	IsVerified    bool      `json:"is_verified"`
	RatingAvg     float64   `json:"rating_avg"`
	TotalSales    int       `json:"total_sales"`
	TotalProducts int       `json:"total_products"`
	CreatedAt     time.Time `json:"created_at"`
	Revenue       float64   `json:"revenue,omitempty"`
}

type AdminPlatformSummary struct {
	TotalUsers         int     `json:"total_users"`
	TotalBuyers        int     `json:"total_buyers"`
	TotalSellers       int     `json:"total_sellers"`
	TotalAdmins        int     `json:"total_admins"`
	TotalStores        int     `json:"total_stores"`
	VerifiedStores     int     `json:"verified_stores"`
	UnverifiedStores   int     `json:"unverified_stores"`
	TotalProducts      int     `json:"total_products"`
	ActiveProducts     int     `json:"active_products"`
	InactiveProducts   int     `json:"inactive_products"`
	TotalOrders        int     `json:"total_orders"`
	PendingOrders      int     `json:"pending_orders"`
	PaidOrders         int     `json:"paid_orders"`
	ProcessingOrders   int     `json:"processing_orders"`
	ShippedOrders      int     `json:"shipped_orders"`
	DeliveredOrders    int     `json:"delivered_orders"`
	CompletedOrders    int     `json:"completed_orders"`
	CanceledOrders     int     `json:"canceled_orders"`
	TotalRevenue       float64 `json:"total_revenue"`
	NewUsersThisWeek   int     `json:"new_users_this_week"`
	NewUsersThisMonth  int     `json:"new_users_this_month"`
	NewOrdersThisWeek  int     `json:"new_orders_this_week"`
	NewOrdersThisMonth int     `json:"new_orders_this_month"`
}

type AdminTopProductResponse struct {
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	StoreName    string  `json:"store_name"`
	CategoryName string  `json:"category_name"`
	Price        float64 `json:"price"`
	RatingAvg    float64 `json:"rating_avg"`
	QuantitySold int     `json:"quantity_sold"`
	Revenue      float64 `json:"revenue"`
}
