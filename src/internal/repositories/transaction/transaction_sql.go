package transaction

import (
	"context"
	"database/sql"

	"ne-project/src/internal/models/dto"
)

func (r *TransactionRepository) GetToday(
	ctx context.Context,
) (*dto.TransactionListResponse, error) {

	var trx dto.TransactionListResponse

	err := r.db.QueryRowContext(ctx, querySummary).
		Scan(&trx.Data)

	if err != nil {
		return nil, err
	}

	// var best dto.BestSellingProduct

	// err = r.db.QueryRowContext(ctx, queryBestProduct).
	// 	Scan(&best.Name, &best.QuantitySold)

	if err == sql.ErrNoRows {

		// summary.BestSelling = dto.BestSellingProduct{
		// 	Name:         "-",
		// 	QuantitySold: 0,
		// }

		return &trx, nil
	}

	// if err != nil {
	// 	return nil, err
	// }

	// summary.BestSelling = best

	return &trx, nil
}
