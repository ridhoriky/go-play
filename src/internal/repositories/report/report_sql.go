package report

import (
	"context"
	"log"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/utils/validation"

	"github.com/rs/zerolog"
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
		zerolog.Ctx(ctx).Error().Err(err).Msg("err query summary")
		return dto.SummaryResponse{}, err
	}
	return summary, nil
}

func (r *reportRepository) GetTopProducts(
	ctx context.Context, dr validation.DateRange, limit int,
) (dto.TopProductsResponse, error) {
	var topProducts dto.TopProductsResponse

	rows, err := r.db.QueryContext(ctx, getTopProductsQuery, dr.From, dr.To, limit)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err query top products")
		return dto.TopProductsResponse{}, err
	}
	defer rows.Close()
	log.Println(rows, "rows")
	for rows.Next() {
		var product dto.TopProductItem
		if err := rows.Scan(&product.ProductID, &product.ProductName, &product.TotalQuantity, &product.TotalRevenue); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("err scanning top products")
			return dto.TopProductsResponse{}, err
		}
		topProducts.Data = append(topProducts.Data, product)
	}
	return topProducts, nil
}
