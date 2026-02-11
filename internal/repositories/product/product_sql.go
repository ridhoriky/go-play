package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ne-project/internal/database"
	"ne-project/internal/dto"
	"ne-project/internal/models"

	"github.com/jmoiron/sqlx"
)

func (repo *ProductRepository) GetAll(ctx context.Context, req dto.ProductFilterRequest) ([]dto.ProductResponse, error){
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

	for rows.Next(){
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

func (repo *ProductRepository) Create(ctx context.Context, product *models.Product) error{
	query := insertProductQuery
	err := repo.db.QueryRowContext(ctx, query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)

	return err
}

func (repo *ProductRepository) Update(ctx context.Context, product *models.Product) error{
	query := updateProductQuery
	result , err := repo.db.ExecContext(ctx, query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
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

func (repo *ProductRepository) Delete(ctx context.Context, id int) error{
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
		return errors.New("Product not found")
	}

	return err
}

func (repo *ProductRepository) GetByID(ctx context.Context, id int) (*dto.ProductResponse, error){
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
	products []models.Product,
) ([]dto.ProductDTO, error) {

	if len(products) == 0 {
		return nil, errors.New("empty products")
	}

	var responses []dto.ProductDTO

	err := database.WithTransactionX(ctx, repo.db, func(tx *sqlx.Tx) error {

		columns := []string{
			"name",
			"price",
			"stock",
			"category_id",
		}

		rowsData := make([][]any, 0, len(products))

		for _, p := range products {

			row := []any{
				p.Name,
				p.Price,
				p.Stock,
				p.CategoryID,
			}

			rowsData = append(rowsData, row)
		}

		query, args, err := database.BuildBulkInsert(
			"products",
			columns,
			rowsData,
		)
		if err != nil {
			return err
		}

		query += `
			RETURNING
				id,
				name,
				price,
				stock,
				category_id
		`

		rows, err := tx.QueryxContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {

			var resp dto.ProductDTO

			if err := rows.StructScan(&resp); err != nil {
				return err
			}

			responses = append(responses, resp)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return responses, nil
}
