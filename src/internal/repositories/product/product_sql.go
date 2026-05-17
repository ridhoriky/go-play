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

	"github.com/rs/zerolog"
)

func (r *productRepository) GetAll(ctx context.Context, filter *dto.GetProductsQuery) ([]entity.ProductWithCategory, int, error) {
	filterQuery, args := buildProductFilters(filter)

	dataQuery := getAllProductsQuery + filterQuery

	sortBy := "p.created_at"
	allowedSortColumns := map[string]string{
		"name":       "p.name",
		"price":      "p.price",
		"stock":      "p.stock",
		"created_at": "p.created_at",
	}

	if val, ok := allowedSortColumns[filter.SortBy]; ok {
		sortBy = val
	}

	sortDir := "ASC"
	if filter.SortDir == "DESC" {
		sortDir = "DESC"
	}

	dataQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortDir)

	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)

	offset := (filter.Page - 1) * filter.Limit

	args = append(args, filter.Limit, offset)

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
			&p.Name,
			&p.Price,
			&p.Stock,
			&p.CategoryID,
			&p.CategoryName,
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

func buildProductFilters(filter *dto.GetProductsQuery) (string, []any) {

	query := ""
	args := []any{}
	argPos := 1

	// search
	if filter.Search != "" {
		query += fmt.Sprintf(" AND p.name ILIKE $%d", argPos)
		args = append(args, "%"+filter.Search+"%")
		argPos++
	}

	// category
	if filter.CategoryID != "" {
		query += fmt.Sprintf(" AND p.category_id = $%d", argPos)
		args = append(args, filter.CategoryID)
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
	}

	return query, args
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	query := insertProductQuery
	err := r.db.QueryRowContext(ctx, query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err create product")
		return err
	}
	return nil
}

func (r *productRepository) Update(ctx context.Context, id string, product *entity.Product) error {
	query := updateProductQuery
	result, err := r.db.ExecContext(ctx, query, product.Name, product.Price, product.Stock, product.CategoryID, id)
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
		return dto.NewError(http.StatusNotFound, preference.ErrProductNotFound)
	}

	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id string) (*entity.Product, string, error) {
	query := getProductByIDQuery

	var p entity.Product
	categoryName := ""

	err := r.db.
		QueryRowContext(ctx, query, id).
		Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Stock,
			&p.CategoryID,
			&categoryName,
		)

	if errors.Is(err, sql.ErrNoRows) {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg(preference.ErrProductNotFound)
		return nil, categoryName, dto.NewError(http.StatusNotFound, preference.ErrProductNotFound)
	}

	if err != nil {
		return nil, categoryName, err
	}

	return &p, categoryName, nil
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
	insertColumns := []string{"name", "price", "stock", "category_id"}

	returningColumns := []string{
		"id",
		"name",
		"price",
		"stock",
		"category_id",
	}

	var (
		values = make([]string, 0, len(products))
		args   = make([]any, 0, 4*len(products))
		argPos = 1
	)

	for i := range products {
		p := &products[i]
		values = append(values,
			fmt.Sprintf("($%d,$%d,$%d,$%d)",
				argPos,
				argPos+1,
				argPos+2,
				argPos+3,
			),
		)

		args = append(args,
			p.Name,
			p.Price,
			p.Stock,
			p.CategoryID,
		)

		argPos += 4
	}

	return fmt.Sprintf(
		insertBulkProductQuery,
		strings.Join(insertColumns, ","),
		strings.Join(values, ","),
		strings.Join(returningColumns, ","),
	), args
}
