package dto

// ─── Request ────────────────────────────────────────────────────────────────

type GetReportQuery struct {
	Period   string `form:"period"    binding:"omitempty,oneof=today this_week this_month custom"`
	DateFrom string `form:"date_from" binding:"omitempty"`
	DateTo   string `form:"date_to"   binding:"omitempty"`
}

type GetTopProductsQuery struct {
	GetReportQuery
	Limit int `form:"limit" binding:"omitempty,min=1,max=20"`
}

// ─── Response ────────────────────────────────────────────────────────────────

type SummaryResponse struct {
	TotalRevenue       float64 `json:"total_revenue"`
	TotalTransactions  int     `json:"total_transactions"`
	TotalItemsSold     int     `json:"total_items_sold"`
	AverageTransaction float64 `json:"average_transaction"`
}

type TopProductItem struct {
	ProductID     string  `json:"product_id"`
	ProductName   string  `json:"product_name"`
	TotalQuantity int     `json:"total_quantity"`
	TotalRevenue  float64 `json:"total_revenue"`
}

type TopProductsResponse struct {
	Data []TopProductItem `json:"data"`
}
