package dto

import "github.com/shopspring/decimal"

type BestSellingProduct struct {
	Name         string `json:"name"`
	QuantitySold int    `json:"quantity_sold"`
}

type TodaySummaryResponse struct {
	TotalRevenue      decimal.Decimal    `json:"total_revenue"`
	TotalTransactions int                `json:"total_transactions"`
	BestSelling       BestSellingProduct `json:"best_selling_product"`
}
