package transaction

import (
	"context"
	"database/sql"

	"ne-project/src/internal/models/dto"
)

func (r *TransactionRepository) GetToday(
	ctx context.Context,
) (*dto.TodaySummaryResponse, error) {

	var summary dto.TodaySummaryResponse

	err := r.db.QueryRowContext(ctx, querySummary).
		Scan(&summary.TotalRevenue, &summary.TotalTransactions)

	if err != nil {
		return nil, err
	}

	var best dto.BestSellingProduct

	err = r.db.QueryRowContext(ctx, queryBestProduct).
		Scan(&best.Name, &best.QuantitySold)

	if err == sql.ErrNoRows {

		summary.BestSelling = dto.BestSellingProduct{
			Name:         "-",
			QuantitySold: 0,
		}

		return &summary, nil
	}

	if err != nil {
		return nil, err
	}

	summary.BestSelling = best

	return &summary, nil
}
