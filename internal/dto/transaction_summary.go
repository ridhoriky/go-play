package dto

type BestSellingProduct struct {
	Name         string `json:"name"`
	QuantitySold int    `json:"quantity_sold"`
}

type TodaySummaryResponse struct {
	TotalRevenue      int                `json:"total_revenue"`
	TotalTransactions int                `json:"total_transactions"`
	BestSelling       BestSellingProduct `json:"best_selling_product"`
}
