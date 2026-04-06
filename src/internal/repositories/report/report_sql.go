package report

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/utils/validation"
)

func (r *reportRepository) GetSummary(
	ctx context.Context, dr validation.DateRange,
) (dto.SummaryResponse, error) {
	summary := dto.SummaryResponse{}
	err := r.db.QueryRowContext(ctx, getSummaryQuery, dr.From, dr.To).
		Scan(
			&summary.TotalRevenue,
			&summary.TotalTransactions,
			&summary.TotalItemsSold,
			&summary.AverageTransaction,
		)
	if err != nil {
		return dto.SummaryResponse{}, err
	}
	return summary, nil
}
