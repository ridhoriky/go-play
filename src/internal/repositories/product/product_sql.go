package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (r *productRepository) GetAll(ctx context.Context, filter *dto.GetProductsQuery) ([]entity.ProductWithCategory, int, error) {
	filterQuery, args := buildProductFilters(filter)

	dataQuery := getAllProductsQuery + filterQuery

	dataQuery += buildProductSort(filter, &args)

	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, filter.Limit, (filter.Page-1)*filter.Limit)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err find product with query")
		return nil, 0, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			zerolog.Ctx(ctx).Error().Err(closeErr).Msg("failed to close rows")
		}
	}()

	products := []entity.ProductWithCategory{}

	var total int
	for rows.Next() {

		var p entity.ProductWithCategory

		if scanErr := rows.Scan(
			&p.ID,
			&p.StoreID,
			&p.CategoryID,
			&p.Name,
			&p.Slug,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.RatingAvg,
			&p.TotalSold,
			&p.IsActive,
			&p.CategoryName,
			&p.StoreName,
			&p.StoreSlug,
			&p.StoreIsVerified,
			&p.PrimaryImage,
			&total,
		); scanErr != nil {
			zerolog.Ctx(ctx).Error().Err(scanErr).Str("productID", p.ID).Msg("err mapping product row")
			return nil, 0, scanErr
		}

		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err iterating product row")
		return nil, 0, err
	}

	return products, total, nil
}

func buildProductSort(filter *dto.GetProductsQuery, args *[]any) string {
	sortBy := "p.created_at"
	allowedSortColumns := map[string]string{
		"newest":     "p.created_at",
		"price_asc":  "p.price",
		"price_desc": "p.price",
		"rating":     "p.rating_avg",
		"popular":    "p.total_sold",
	}

	if val, ok := allowedSortColumns[filter.Sort]; ok {
		sortBy = val
	}

	sortDir := "ASC"
	if filter.Sort == "newest" || filter.Sort == "price_desc" || filter.Sort == "rating" || filter.Sort == "popular" {
		sortDir = "DESC"
	}

	if filter.Sort == "" && filter.Q != "" {
		sortBy = fmt.Sprintf("ts_rank(p.search_vector, plainto_tsquery('english', $%d)) + CASE WHEN s.is_verified THEN 0.2 ELSE 0.0 END", len(*args)+1)
		*args = append(*args, filter.Q)
		sortDir = "DESC"
	}

	return fmt.Sprintf(" ORDER BY %s %s", sortBy, sortDir)
}

func buildProductFilters(filter *dto.GetProductsQuery) (string, []any) {

	query := ""
	args := []any{}
	argPos := 1

	// search
	if filter.Q != "" {
		query += fmt.Sprintf(" AND p.search_vector @@ plainto_tsquery('english', $%d)", argPos)
		args = append(args, filter.Q)
		argPos++
	}

	// category
	if filter.Category != "" {
		if _, err := uuid.Parse(filter.Category); err == nil {
			query += fmt.Sprintf(" AND p.category_id = $%d", argPos)
		} else {
			query += fmt.Sprintf(" AND (c.name ILIKE $%d OR LOWER(c.name) = LOWER($%d) OR REPLACE(LOWER(c.name), ' & ', '-') = LOWER($%d) OR REPLACE(LOWER(c.name), ' ', '-') = LOWER($%d) OR LOWER(c.name) ILIKE '%%' || LOWER($%d) || '%%')", argPos, argPos, argPos, argPos, argPos)
		}
		args = append(args, filter.Category)
		argPos++
	}

	// store
	if filter.Store != "" {
		query += fmt.Sprintf(" AND p.store_id = $%d", argPos)
		args = append(args, filter.Store)
		argPos++
	}

	// min price
	if !filter.MinPrice.IsZero() {
		query += fmt.Sprintf(" AND p.price >= $%d", argPos)
		args = append(args, filter.MinPrice)
		argPos++
	}

	// max price
	if !filter.MaxPrice.IsZero() {
		query += fmt.Sprintf(" AND p.price <= $%d", argPos)
		args = append(args, filter.MaxPrice)
		argPos++
	}

	// rating
	if filter.Rating > 0 {
		query += fmt.Sprintf(" AND p.rating_avg >= $%d", argPos)
		args = append(args, filter.Rating)
		argPos++
	}

	// stock filter
	if filter.InStock {
		query += " AND p.stock > 0"
	}

	if filter.LowStock {
		query += " AND p.stock <= 5 AND p.stock > 0"
	}

	if filter.IsActive != nil {
		query += fmt.Sprintf(" AND p.is_active = $%d", argPos)
		args = append(args, *filter.IsActive)
	}

	return query, args
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	query := insertProductQuery
	err := r.db.QueryRowContext(ctx, query, product.StoreID, product.CategoryID, product.Name, product.Slug, product.Description, product.Price, product.Stock, product.IsActive).Scan(&product.ID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err create product")
		return err
	}
	return nil
}

func (r *productRepository) Update(ctx context.Context, id string, product *entity.Product) error {
	query := updateProductQuery
	result, err := r.db.ExecContext(ctx, query, product.CategoryID, product.Name, product.Description, product.Price, product.Stock, product.IsActive, id, product.StoreID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err update product")
		return err
	}

	rowCount, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err rows affected")
		return err
	}

	if rowCount == 0 {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg(preference.ErrProductNotFound)
		return dto.NewError(http.StatusNotFound, preference.ErrProductNotFound)
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	query := deleteProductQuery
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err delete product")
		return err
	}
	rowCount, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err rows affected")
		return err
	}

	if rowCount == 0 {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg(preference.ErrProductNotFound)
		return dto.NewError(http.StatusNotFound, preference.ErrProductNotFound)
	}

	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*entity.ProductDetail, error) {
	query := getProductByIDQuery

	var p entity.ProductDetail

	err := r.db.
		QueryRowContext(ctx, query, id).
		Scan(
			&p.ID,
			&p.StoreID,
			&p.CategoryID,
			&p.Name,
			&p.Slug,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.RatingAvg,
			&p.TotalSold,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.CategoryName,
			&p.StoreName,
			&p.StoreSlug,
			&p.StoreIsVerified,
			&p.StoreLogoURL,
			&p.StoreRatingAvg,
			&p.TotalReviews,
		)

	if errors.Is(err, sql.ErrNoRows) {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg(preference.ErrProductNotFound)
		return nil, dto.NewError(http.StatusNotFound, preference.ErrProductNotFound)
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *productRepository) GetBySlug(ctx context.Context, slug string) (*entity.Product, error) {
	query := getProductBySlugQuery

	var p entity.Product

	err := r.db.
		QueryRowContext(ctx, query, slug).
		Scan(
			&p.ID,
			&p.StoreID,
			&p.CategoryID,
			&p.Name,
			&p.Slug,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.RatingAvg,
			&p.TotalSold,
			&p.IsActive,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, dto.NewError(http.StatusNotFound, preference.ErrProductNotFound)
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *productRepository) GetDetailBySlug(ctx context.Context, slug string) (*entity.ProductDetail, error) {
	query := getProductDetailBySlugQuery

	var p entity.ProductDetail

	err := r.db.
		QueryRowContext(ctx, query, slug).
		Scan(
			&p.ID,
			&p.StoreID,
			&p.CategoryID,
			&p.Name,
			&p.Slug,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.RatingAvg,
			&p.TotalSold,
			&p.IsActive,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.CategoryName,
			&p.StoreName,
			&p.StoreSlug,
			&p.StoreIsVerified,
			&p.StoreLogoURL,
			&p.StoreRatingAvg,
			&p.TotalReviews,
		)

	if errors.Is(err, sql.ErrNoRows) {
		zerolog.Ctx(ctx).Error().Err(err).Str("slug", slug).Msg(preference.ErrProductNotFound)
		return nil, dto.NewError(http.StatusNotFound, preference.ErrProductNotFound)
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *productRepository) CreateMultiple(
	ctx context.Context,
	products []entity.Product,
) ([]entity.Product, error) {

	if len(products) == 0 {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrProductEmpty)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrProductEmpty)
	}

	if len(products) > preference.MaxBatchSizeProduct {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrProductBatchTooLarge)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrProductBatchTooLarge)
	}

	timeout := min(
		time.Duration(len(products))*200*time.Millisecond,
		30*time.Second,
	)

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})

	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err tx create product")
		return nil, err
	}

	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			zerolog.Ctx(ctx).Error().Err(rbErr).Msg("err rollback create product")
		}
	}()

	query, args := buildBulkInsertQuery(products)

	rows, err := tx.QueryxContext(ctx, query, args...)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err query create products")
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			zerolog.Ctx(ctx).Error().Err(closeErr).Msg("failed to close rows")
		}
	}()

	var responses []entity.Product

	for rows.Next() {

		var resp entity.Product

		if scanErr := rows.StructScan(&resp); scanErr != nil {
			zerolog.Ctx(ctx).Error().Err(scanErr).Interface("row_data", resp).Msg("failed to scan product row")
			return nil, scanErr
		}

		responses = append(responses, resp)
	}

	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err iterate product row")
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err commit create products")
		return nil, err
	}

	return responses, nil
}

func buildBulkInsertQuery(products []entity.Product) (string, []any) {
	insertColumns := []string{"store_id", "category_id", "name", "slug", "description", "price", "stock", "is_active"}

	returningColumns := []string{
		"id",
		"store_id",
		"category_id",
		"name",
		"slug",
		"description",
		"price",
		"stock",
		"is_active",
	}

	var (
		values = make([]string, 0, len(products))
		args   = make([]any, 0, 8*len(products))
		argPos = 1
	)

	for i := range products {
		p := &products[i]
		values = append(values,
			fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
				argPos,
				argPos+1,
				argPos+2,
				argPos+3,
				argPos+4,
				argPos+5,
				argPos+6,
				argPos+7,
			),
		)

		args = append(args,
			p.StoreID,
			p.CategoryID,
			p.Name,
			p.Slug,
			p.Description,
			p.Price,
			p.Stock,
			p.IsActive,
		)

		argPos += 8
	}

	return fmt.Sprintf(
		insertBulkProductQuery,
		strings.Join(insertColumns, ","),
		strings.Join(values, ","),
		strings.Join(returningColumns, ","),
	), args
}
