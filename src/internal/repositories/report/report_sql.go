package report

import (
	"context"

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
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			zerolog.Ctx(ctx).Error().Err(closeErr).Msg("failed to close rows")
		}
	}()
	for rows.Next() {
		var product dto.TopProductItem
		if err = rows.Scan(&product.ProductID, &product.ProductName, &product.TotalQuantity, &product.TotalRevenue); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("err scanning top products")
			return dto.TopProductsResponse{}, err
		}
		topProducts.Data = append(topProducts.Data, product)
	}
	return topProducts, nil
}

func (r *reportRepository) GetSellerSalesSummary(
	ctx context.Context, storeID string, dr validation.DateRange,
) (*dto.SellerSalesSummary, error) {
	summary := &dto.SellerSalesSummary{}
	err := r.db.QueryRowContext(ctx, getSellerSalesSummaryQuery, storeID, dr.From, dr.To).
		Scan(
			&summary.TotalOrders,
			&summary.TotalRevenue,
			&summary.TotalItemsSold,
		)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err query seller sales summary")
		return nil, err
	}
	return summary, nil
}

func (r *reportRepository) GetSellerTopProducts(
	ctx context.Context, storeID string, dr validation.DateRange, limit int, offset int,
) ([]dto.SellerTopProduct, error) {
	rows, err := r.db.QueryContext(ctx, getSellerTopProductsQuery, storeID, dr.From, dr.To, limit, offset)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err query seller top products")
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			zerolog.Ctx(ctx).Error().Err(closeErr).Msg("failed to close rows")
		}
	}()

	var products []dto.SellerTopProduct
	for rows.Next() {
		var p dto.SellerTopProduct
		if err := rows.Scan(&p.ProductID, &p.ProductName, &p.TotalSold, &p.Revenue, &p.AverageRating); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("err scan seller top products")
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *reportRepository) GetSellerRecentOrders(
	ctx context.Context, storeID string, limit int,
) ([]dto.OrderSummary, error) {
	rows, err := r.db.QueryContext(ctx, getSellerRecentOrdersQuery, storeID, limit)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err query seller recent orders")
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			zerolog.Ctx(ctx).Error().Err(closeErr).Msg("failed to close rows")
		}
	}()

	var orders []dto.OrderSummary
	for rows.Next() {
		var o dto.OrderSummary
		if err := rows.Scan(&o.ID, &o.TotalAmount, &o.Status, &o.CreatedAt, &o.BuyerName); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("err scan seller recent orders")
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func (r *reportRepository) GetSellerPendingOrdersCount(
	ctx context.Context, storeID string,
) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, getSellerPendingOrdersCountQuery, storeID).Scan(&count)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err query seller pending orders count")
		return 0, err
	}
	return count, nil
}

func (r *reportRepository) GetSellerLowStockProducts(
	ctx context.Context, storeID string, threshold int,
) ([]dto.ProductLowStock, error) {
	rows, err := r.db.QueryContext(ctx, getSellerLowStockProductsQuery, storeID, threshold)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err query seller low stock products")
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			zerolog.Ctx(ctx).Error().Err(closeErr).Msg("failed to close rows")
		}
	}()

	var products []dto.ProductLowStock
	for rows.Next() {
		var p dto.ProductLowStock
		if err := rows.Scan(&p.ID, &p.Name, &p.Stock); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("err scan seller low stock products")
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
