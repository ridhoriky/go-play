package report

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/utils/validation"

	"github.com/rs/zerolog"
)

func (s *reportService) GetSummary(
	ctx context.Context, req *dto.GetReportQuery,
) (*dto.SummaryResponse, error) {
	dr, err := validation.Parse(req.Period, req.DateFrom, req.DateTo)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg(preference.ErrInvalidQueryParams)
		return nil, err
	}

	summary, err := s.reportRepository.GetSummary(ctx, dr)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg(preference.ErrInvalidQueryParams)
		return nil, err
	}

	return &dto.SummaryResponse{
		TotalRevenue:       summary.TotalRevenue,
		TotalTransactions:  summary.TotalTransactions,
		TotalItemsSold:     summary.TotalItemsSold,
		AverageTransaction: summary.AverageTransaction,
	}, nil
}

func (s *reportService) GetTopProducts(
	ctx context.Context, req *dto.GetTopProductsQuery,
) (*dto.TopProductsResponse, error) {
	dr, err := validation.Parse(req.Period, req.DateFrom, req.DateTo)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg(preference.ErrInvalidQueryParams)
		return nil, err
	}

	req.Limit = validation.ValidatePageSize(req.Limit)

	topProducts, err := s.reportRepository.GetTopProducts(ctx, dr, req.Limit)
	if err != nil {
		return nil, err
	}

	return &dto.TopProductsResponse{
		Data: topProducts.Data,
	}, nil
}
