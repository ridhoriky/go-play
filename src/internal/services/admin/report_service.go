package admin

import (
	"context"

	"ne-project/src/internal/models/dto"
	repository "ne-project/src/internal/repositories/admin"
)

type AdminReportServiceItf interface {
	GetPlatformSummary(ctx context.Context) (*dto.AdminPlatformSummary, error)
	GetTopStores(ctx context.Context, sortBy string, limit int) ([]dto.AdminSellerResponse, error)
	GetTopProducts(ctx context.Context, sortBy string, limit int) ([]dto.AdminTopProductResponse, error)
}

type adminReportService struct {
	repo repository.AdminRepositoryItf
}

func NewAdminReportService(repo repository.AdminRepositoryItf) AdminReportServiceItf {
	return &adminReportService{
		repo: repo,
	}
}

func (s *adminReportService) GetPlatformSummary(ctx context.Context) (*dto.AdminPlatformSummary, error) {
	return s.repo.GetPlatformSummary(ctx)
}

func (s *adminReportService) GetTopStores(ctx context.Context, sortBy string, limit int) ([]dto.AdminSellerResponse, error) {
	validSortBy := map[string]bool{
		"revenue":     true,
		"total_sales": true,
		"rating":      true,
		"products":    true,
	}

	if !validSortBy[sortBy] {
		sortBy = "revenue"
	}

	if limit <= 0 || limit > 50 {
		limit = 10
	}

	return s.repo.GetTopStores(ctx, sortBy, limit)
}

func (s *adminReportService) GetTopProducts(ctx context.Context, sortBy string, limit int) ([]dto.AdminTopProductResponse, error) {
	validSortBy := map[string]bool{
		"revenue":  true,
		"quantity": true,
		"rating":   true,
	}

	if !validSortBy[sortBy] {
		sortBy = "revenue"
	}

	if limit <= 0 || limit > 50 {
		limit = 10
	}

	return s.repo.GetTopProducts(ctx, sortBy, limit)
}
