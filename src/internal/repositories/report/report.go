package report

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/utils/validation"

	"github.com/jmoiron/sqlx"
)

type ReportRepositoryItf interface {
	GetSummary(ctx context.Context, dr validation.DateRange) (dto.SummaryResponse, error)
	GetTopProducts(ctx context.Context, dr validation.DateRange, limit int) (dto.TopProductsResponse, error)

	GetSellerSalesSummary(ctx context.Context, storeID string, dr validation.DateRange) (*dto.SellerSalesSummary, error)
	GetSellerTopProducts(ctx context.Context, storeID string, dr validation.DateRange, limit int, offset int) ([]dto.SellerTopProduct, error)
	GetSellerRecentOrders(ctx context.Context, storeID string, limit int) ([]dto.OrderSummary, error)
	GetSellerPendingOrdersCount(ctx context.Context, storeID string) (int, error)
	GetSellerLowStockProducts(ctx context.Context, storeID string, threshold int) ([]dto.ProductLowStock, error)
}
type reportRepository struct {
	db *sqlx.DB
}

func NewReportRepository(db *sqlx.DB) ReportRepositoryItf {
	return &reportRepository{
		db: db,
	}
}
