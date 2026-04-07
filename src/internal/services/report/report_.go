package report

import (
	"context"
	"log"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/utils/validation"
)

func (s *reportService) GetSummary(
	ctx context.Context, req *dto.GetReportQuery,
) (*dto.SummaryResponse, error) {
	log.Println(req, "req")
	dr, err := validation.Parse(req.Period, req.DateFrom, req.DateTo)
	if err != nil {
		return nil, err
	}
	log.Println(dr, "req dr")

	summary, err := s.reportRepository.GetSummary(ctx, dr)
	if err != nil {
		return nil, err
	}

	log.Println("summary", summary)

	return &dto.SummaryResponse{
		TotalRevenue:       summary.TotalRevenue,
		TotalTransactions:  summary.TotalTransactions,
		TotalItemsSold:     summary.TotalItemsSold,
		AverageTransaction: summary.AverageTransaction,
	}, nil
}
