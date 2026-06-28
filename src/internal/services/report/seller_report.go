package report

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/repositories/report"
	"ne-project/src/internal/utils/validation"

	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

type SellerReportServiceItf interface {
	GetSalesSummary(ctx context.Context, storeID string, req *dto.GetSellerReportQuery) (*dto.SellerSalesSummary, error)
	GetTopProducts(ctx context.Context, storeID string, req *dto.GetSellerReportQuery, limit int, offset int) ([]dto.SellerTopProduct, error)
	GetDashboard(ctx context.Context, storeID string, req *dto.GetSellerReportQuery) (*dto.SellerDashboard, error)
}

type sellerReportService struct {
	reportRepo report.ReportRepositoryItf
}

func NewSellerReportService(reportRepo report.ReportRepositoryItf) SellerReportServiceItf {
	return &sellerReportService{
		reportRepo: reportRepo,
	}
}

func (s *sellerReportService) GetSalesSummary(
	ctx context.Context, storeID string, req *dto.GetSellerReportQuery,
) (*dto.SellerSalesSummary, error) {
	if req.Period == "" {
		req.Period = "this_week"
	}
	dr, err := validation.Parse(req.Period, req.DateFrom, req.DateTo)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg(preference.ErrInvalidQueryParams)
		return nil, err
	}

	summary, err := s.reportRepo.GetSellerSalesSummary(ctx, storeID, dr)
	if err != nil {
		return nil, err
	}
	summary.Period = dr.From.Format("2006-01-02") + " to " + dr.To.Format("2006-01-02")

	if summary.TotalOrders > 0 {
		summary.AverageOrderValue = summary.TotalRevenue.Div(decimal.NewFromInt(int64(summary.TotalOrders)))
	}

	return summary, nil
}

func (s *sellerReportService) GetTopProducts(
	ctx context.Context, storeID string, req *dto.GetSellerReportQuery, limit int, offset int,
) ([]dto.SellerTopProduct, error) {
	if req.Period == "" {
		req.Period = "this_week"
	}
	dr, err := validation.Parse(req.Period, req.DateFrom, req.DateTo)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg(preference.ErrInvalidQueryParams)
		return nil, err
	}

	if limit <= 0 {
		limit = 10
	}

	products, err := s.reportRepo.GetSellerTopProducts(ctx, storeID, dr, limit, offset)
	if err != nil {
		return nil, err
	}
	if products == nil {
		products = []dto.SellerTopProduct{}
	}

	return products, nil
}

func (s *sellerReportService) GetDashboard(
	ctx context.Context, storeID string, req *dto.GetSellerReportQuery,
) (*dto.SellerDashboard, error) {
	if req.Period == "" {
		req.Period = "this_week"
	}
	dr, err := validation.Parse(req.Period, req.DateFrom, req.DateTo)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg(preference.ErrInvalidQueryParams)
		return nil, err
	}

	// Sales Summary
	summary, err := s.reportRepo.GetSellerSalesSummary(ctx, storeID, dr)
	if err != nil {
		return nil, err
	}
	summary.Period = dr.From.Format("2006-01-02") + " to " + dr.To.Format("2006-01-02")
	if summary.TotalOrders > 0 {
		summary.AverageOrderValue = summary.TotalRevenue.Div(decimal.NewFromInt(int64(summary.TotalOrders)))
	}

	// Top Products (limit 5)
	topProducts, err := s.reportRepo.GetSellerTopProducts(ctx, storeID, dr, 5, 0)
	if err != nil {
		return nil, err
	}
	if topProducts == nil {
		topProducts = []dto.SellerTopProduct{}
	}

	// Recent Orders (limit 5)
	recentOrders, err := s.reportRepo.GetSellerRecentOrders(ctx, storeID, 5)
	if err != nil {
		return nil, err
	}
	if recentOrders == nil {
		recentOrders = []dto.OrderSummary{}
	}

	// Pending Orders Count
	pendingCount, err := s.reportRepo.GetSellerPendingOrdersCount(ctx, storeID)
	if err != nil {
		return nil, err
	}

	// Low Stock Items (threshold = 10)
	lowStock, err := s.reportRepo.GetSellerLowStockProducts(ctx, storeID, 10)
	if err != nil {
		return nil, err
	}
	if lowStock == nil {
		lowStock = []dto.ProductLowStock{}
	}

	return &dto.SellerDashboard{
		Summary:       *summary,
		TopProducts:   topProducts,
		RecentOrders:  recentOrders,
		PendingOrders: pendingCount,
		LowStockItems: lowStock,
	}, nil
}
