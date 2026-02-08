package product

import (
	"context"
	"database/sql"
	"errors"
	"ne-project/internal/database"
	"ne-project/internal/dto"
	"ne-project/internal/models"
)

func (repo *ProductRepository) GetAll(ctx context.Context) ([]dto.ProductResponse, error){
	query := getAllProductsQuery

	rows, err := repo.db.QueryContext(ctx, query)

	if err != nil {
		return nil , err 
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

func (repo *ProductRepository) CreateMultiple(ctx context.Context, products []models.Product) ([]dto.ProductResponse, error) {

	var responses []dto.ProductResponse

	err := database.WithTransaction(ctx, repo.db, func(tx *sql.Tx) error {

		insertStmt, err := tx.PrepareContext(ctx, insertProductQuery)
		if err != nil {
			return err
		}
		defer insertStmt.Close()

		selectStmt, err := tx.PrepareContext(ctx, getProductByIDQuery)
		if err != nil {
			return err
		}
		defer selectStmt.Close()

		responses = make([]dto.ProductResponse, 0, len(products))

		for _, p := range products {

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			var productID int

			err := insertStmt.QueryRowContext(
				ctx,
				p.Name,
				p.Price,
				p.Stock,
				p.CategoryID,
			).Scan(&productID)

			if err != nil {
				return err
			}

			var resp dto.ProductResponse

			err = selectStmt.QueryRowContext(ctx, productID).
				Scan(
					&resp.ID,
					&resp.Name,
					&resp.Price,
					&resp.Stock,
					&resp.CategoryID,
					&resp.CategoryName,
					&resp.CategoryDescription,
				)

			if err != nil {
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