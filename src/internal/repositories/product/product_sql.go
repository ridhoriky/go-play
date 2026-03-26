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
)

func (repo *ProductRepository) GetAll(ctx context.Context, filter *dto.GetProductsQuery) ([]entity.ProductWithCategory, int, error) {
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
		return nil, 0, err
	}

	defer rows.Close()

	products := []entity.ProductWithCategory{}

	for rows.Next() {

		var p entity.ProductWithCategory
		var categoryName string

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Stock,
			&p.CategoryID,
			&categoryName,
		)

		if err != nil {
			return nil, 0, err
		}

		products = append(products, p)
	}

	countQuery := countProductsQuery + filterQuery

	var total int

	err = repo.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
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

func (repo *ProductRepository) Create(ctx context.Context, product *entity.Product) error {
	query := insertProductQuery
	err := repo.db.QueryRowContext(ctx, query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)

	return err
}

func (repo *ProductRepository) Update(ctx context.Context, id string, product *entity.Product) error {
	query := updateProductQuery
	result, err := repo.db.ExecContext(ctx, query, product.Name, product.Price, product.Stock, product.CategoryID, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return &dto.Error{
			Code:    http.StatusNotFound,
			Message: preference.ErrProductNotFound,
		}
	}

	return nil
}

func (repo *ProductRepository) Delete(ctx context.Context, id string) error {
	query := deleteProductQuery
	result, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return &dto.Error{
			Code:    http.StatusNotFound,
			Message: preference.ErrProductNotFound,
		}
	}

	return err
}

func (repo *ProductRepository) GetByID(ctx context.Context, id string) (*entity.Product, string, error) {
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
		return nil, categoryName, &dto.Error{
			Code:    http.StatusNotFound,
			Message: preference.ErrProductNotFound,
		}
	}

	if err != nil {
		return nil, categoryName, err
	}

	return &p, categoryName, nil
}

func (repo *ProductRepository) CreateMultiple(
	ctx context.Context,
	products []entity.Product,
) ([]entity.Product, error) {

	if len(products) == 0 {
		return nil, &dto.Error{
			Code:    http.StatusInternalServerError,
			Message: preference.ErrProductEmpty,
		}
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var err error

	tx, err := repo.db.BeginTxx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
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
		return nil, err
	}
	defer rows.Close()

	var responses []entity.Product

	for rows.Next() {

		var resp entity.Product

		if err := rows.StructScan(&resp); err != nil {
			return nil, err
		}

		responses = append(responses, resp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return responses, nil
}
