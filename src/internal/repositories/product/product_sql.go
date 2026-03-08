package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
)

func (repo *ProductRepository) GetAll(ctx context.Context, req dto.ProductFilterRequest) ([]dto.ProductResponse, error) {
	query := getAllProductsQuery

	args := []any{}
	i := 1

	if req.Name != "" {
		query += fmt.Sprintf(" AND p.name ILIKE $%d", i)
		args = append(args, "%"+req.Name+"%")
		i++
	}

	if req.Category != "" {
		query += fmt.Sprintf(" AND c.name ILIKE $%d", i)
		args = append(args, "%"+req.Category+"%")
		i++
	}
	query += " ORDER BY p.id"

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, req.Limit, (req.Page-1)*req.Limit)

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := make([]dto.ProductResponse, 0)

	for rows.Next() {
		var p dto.ProductResponse
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName, &p.CategoryDescription)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (repo *ProductRepository) Create(ctx context.Context, product *entity.Product) error {
	query := insertProductQuery
	err := repo.db.QueryRowContext(ctx, query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)

	return err
}

func (repo *ProductRepository) Update(ctx context.Context, product *entity.Product) error {
	query := updateProductQuery
	result, err := repo.db.ExecContext(ctx, query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Product not found")
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
		return errors.New(`Product + MESSAGE_404`)
	}

	return err
}

func (repo *ProductRepository) GetByID(ctx context.Context, id string) (*dto.ProductResponse, error) {
	query := getProductByIDQuery

	var p dto.ProductResponse

	err := repo.db.
		QueryRowContext(ctx, query, id).
		Scan(
			&p.ID,
			&p.Name,
			&p.Price,
			&p.Stock,
			&p.CategoryID,
			&p.CategoryName,
			&p.CategoryDescription,
		)

	if err == sql.ErrNoRows {
		return nil, errors.New("Product not found")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) CreateMultiple(
	ctx context.Context,
	products []entity.Product,
) ([]dto.ProductDTO, error) {

	if len(products) == 0 {
		return nil, errors.New("empty products")
	}

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	committed := false

	defer func() {
		if !committed {
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

	var responses []dto.ProductDTO

	for rows.Next() {

		var resp dto.ProductDTO

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

	committed = true

	return responses, nil
}
