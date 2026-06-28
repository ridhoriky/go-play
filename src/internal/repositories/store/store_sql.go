package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"

	"github.com/rs/zerolog"
)

func (r *storeRepository) Create(ctx context.Context, store *entity.Store) error {
	_, err := r.db.ExecContext(ctx, createStoreQuery,
		store.ID, store.UserID, store.StoreName, store.Slug,
		store.Description, store.LogoURL, store.BannerURL,
		store.IsVerified, store.RatingAvg, store.TotalSales,
		store.CreatedAt, store.UpdatedAt,
	)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err to create store")
		return err
	}
	return nil
}

func (r *storeRepository) GetByID(ctx context.Context, id string) (*entity.Store, error) {
	var s entity.Store
	err := r.db.QueryRowContext(ctx, getStoreByIDQuery, id).Scan(
		&s.ID, &s.UserID, &s.StoreName, &s.Slug, &s.Description,
		&s.LogoURL, &s.BannerURL, &s.IsVerified, &s.RatingAvg,
		&s.TotalSales, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err get store by id")
		return nil, err
	}
	return &s, nil
}

func (r *storeRepository) GetByUserID(ctx context.Context, userID string) (*entity.Store, error) {
	var s entity.Store
	err := r.db.QueryRowContext(ctx, getStoreByUserIDQuery, userID).Scan(
		&s.ID, &s.UserID, &s.StoreName, &s.Slug, &s.Description,
		&s.LogoURL, &s.BannerURL, &s.IsVerified, &s.RatingAvg,
		&s.TotalSales, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("userID", userID).Msg("err get store by user id")
		return nil, err
	}
	return &s, nil
}

func (r *storeRepository) GetBySlug(ctx context.Context, slug string) (*entity.Store, error) {
	var s entity.Store
	err := r.db.QueryRowContext(ctx, getStoreBySlugQuery, slug).Scan(
		&s.ID, &s.UserID, &s.StoreName, &s.Slug, &s.Description,
		&s.LogoURL, &s.BannerURL, &s.IsVerified, &s.RatingAvg,
		&s.TotalSales, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("slug", slug).Msg("err get store by slug")
		return nil, err
	}
	return &s, nil
}

func (r *storeRepository) Update(ctx context.Context, store *entity.Store) error {
	result, err := r.db.ExecContext(ctx, updateStoreQuery,
		store.StoreName, store.Description, store.LogoURL, store.BannerURL, store.ID,
	)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", store.ID).Msg("err update store")
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", store.ID).Msg("err get rows affected update store")
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *storeRepository) SoftDelete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, softDeleteStoreQuery, id)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err delete store")
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err get rows affected delete store")
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *storeRepository) IsSlugExists(ctx context.Context, slug string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, checkSlugExistsQuery, slug).Scan(&exists)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("slug", slug).Msg("err check slug exists")
		return false, err
	}
	return exists, nil
}

func (r *storeRepository) List(ctx context.Context, params *dto.GetStoresQuery) ([]entity.Store, int, error) {
	query := listStoresBaseQuery
	args := []any{}

	if params.Search != "" {
		query += " AND store_name ILIKE $1"
		args = append(args, "%"+params.Search+"%")
	}

	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)

	offset := (params.Page - 1) * params.Limit
	args = append(args, params.Limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err find stores")
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	stores := []entity.Store{}
	var total int

	for rows.Next() {
		var s entity.Store
		var rowCount int

		if scanErr := rows.Scan(
			&s.ID, &s.UserID, &s.StoreName, &s.Slug, &s.Description,
			&s.LogoURL, &s.BannerURL, &s.IsVerified, &s.RatingAvg,
			&s.TotalSales, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &rowCount,
		); scanErr != nil {
			zerolog.Ctx(ctx).Error().Err(scanErr).Msg("err mapping store row")
			return nil, 0, scanErr
		}

		stores = append(stores, s)
		if total == 0 {
			total = rowCount
		}
	}

	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err iterating store rows")
		return nil, 0, err
	}

	return stores, total, nil
}

func (r *storeRepository) GetStoreStats(ctx context.Context, storeID string) (*dto.StoreStats, error) {
	var stats dto.StoreStats
	err := r.db.QueryRowContext(ctx, getStoreStatsQuery, storeID).Scan(
		&stats.TotalProducts,
		&stats.ActiveProducts,
		&stats.TotalOrders,
		&stats.TotalRevenue,
		&stats.AverageRating,
		&stats.TotalReviews,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("storeID", storeID).Msg("err get store stats")
		return nil, err
	}
	return &stats, nil
}
