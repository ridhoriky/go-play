package product

import (
	"context"
	"database/sql"
	"errors"
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
	query := createProductQuery
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