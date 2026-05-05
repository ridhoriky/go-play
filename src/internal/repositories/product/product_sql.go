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

func (repo *productRepository) GetAll(ctx context.Context, filter *dto.GetProductsQuery) ([]entity.ProductWithCategory, int, error) {
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

	argsData := append(args, filter.Limit, offset)

	rows, err := repo.db.QueryContext(ctx, dataQuery, argsData...)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err find product with query")
		return nil, 0, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to close rows")
		}
	}()

	products := []entity.ProductWithCategory{}

	var total int
	for rows.Next() {

		var p entity.ProductWithCategory

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Stock,
			&p.CategoryID,
			&p.CategoryName,
			&total,
		)

		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Str("productID", fmt.Sprintf("%v", p.ID)).Msg("err mapping product row")
			return nil, 0, err
		}

		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err iterating product row")
		return nil, 0, err
	}

	return products, total, nil
}

func buildProductFilters(filter *dto.GetProductsQuery) (string, []interface{}) {

	query := ""
	args := []interface{}{}
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
		argPos++
	}

	// stock filter
	if filter.InStock {
		query += " AND p.stock > 0"
	}

	return query, args
}

func (repo *productRepository) Create(ctx context.Context, product *entity.Product) error {
	query := insertProductQuery
	err := repo.db.QueryRowContext(ctx, query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err create product")
		return err
	}
	return nil
}

func (repo *productRepository) Update(ctx context.Context, id string, product *entity.Product) error {
	query := updateProductQuery
	result, err := repo.db.ExecContext(ctx, query, product.Name, product.Price, product.Stock, product.CategoryID, id)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err update product")
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err rows affected")
		return err
	}

	if rows == 0 {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg(preference.ErrProductNotFound)
		return dto.NewError(http.StatusNotFound, preference.ErrProductNotFound)
	}

	return nil
}

func (repo *productRepository) Delete(ctx context.Context, id string) error {
	query := deleteProductQuery
	result, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err delete product")
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err rows affected")
		return err
	}

	if rows == 0 {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg(preference.ErrProductNotFound)
		return dto.NewError(http.StatusNotFound, preference.ErrProductNotFound)
	}

	return err
}

func (repo *productRepository) GetByID(ctx context.Context, id string) (*entity.Product, string, error) {
	query := getProductByIDQuery

	var p entity.Product
	categoryName := ""

	err := repo.db.
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

func (repo *productRepository) CreateMultiple(
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

	var err error

	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})

	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err tx create product")
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			zerolog.Ctx(ctx).Error().Err(err).Msg("err rollback create product")
		}
	}()

	insertColumns := []string{"name", "price", "stock", "category_id"}

	returningColumns := []string{
		"id",
		"name",
		"price",
		"stock",
		"category_id",
	}

	var (
		values []string
		args   []any
		argPos = 1
	)

	for _, p := range products {

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

	query := fmt.Sprintf(
		insertBulkProductQuery,
		strings.Join(insertColumns, ","),
		strings.Join(values, ","),
		strings.Join(returningColumns, ","),
	)

	rows, err := tx.QueryxContext(ctx, query, args...)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err query create products")
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to close rows")
		}
	}()

	var responses []entity.Product

	for rows.Next() {

		var resp entity.Product

		if err := rows.StructScan(&resp); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Interface("row_data", resp).Msg("failed to scan product row")
			return nil, err
		}

		responses = append(responses, resp)
	}

	if err := rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err iterate product row")
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err commit create products")
		return nil, err
	}

	return responses, nil
}
