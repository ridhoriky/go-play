package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type SellerSalesSummary struct {
	Period            string          `json:"period"`
	TotalRevenue      decimal.Decimal `json:"total_revenue"`
	TotalOrders       int             `json:"total_orders"`
	TotalItemsSold    int             `json:"total_items_sold"`
	AverageOrderValue decimal.Decimal `json:"average_order_value"`
	ConversionRate    float64         `json:"conversion_rate,omitempty"`
}

type SellerTopProduct struct {
	ProductID     string          `json:"product_id"`
	ProductName   string          `json:"product_name"`
	TotalSold     int             `json:"total_sold"`
	Revenue       decimal.Decimal `json:"revenue"`
	AverageRating float64         `json:"average_rating"`
}

type OrderSummary struct {
	ID          string          `json:"id"`
	BuyerName   string          `json:"buyer_name"`
	TotalAmount decimal.Decimal `json:"total_amount"`
	Status      string          `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
}

type ProductLowStock struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

type SellerDashboard struct {
	Summary       SellerSalesSummary `json:"summary"`
	TopProducts   []SellerTopProduct `json:"top_products"`
	RecentOrders  []OrderSummary     `json:"recent_orders"`
	PendingOrders int                `json:"pending_orders"`
	LowStockItems []ProductLowStock  `json:"low_stock_items"`
}

type GetSellerReportQuery struct {
	Period   string `form:"period,default=this_week"` // today, this_week, this_month, custom
	DateFrom string `form:"date_from"`
	DateTo   string `form:"date_to"`
}
