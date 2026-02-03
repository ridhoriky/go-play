package product

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"ne-project/internal/dto"
	"ne-project/internal/models"
)

func (repo *ProductRepository) GetAll(ctx context.Context) ([]dto.ProductResponse, error){
	query := `SELECT p.id, p.name, p.price, p.stock, p.category_id, 
			COALESCE(c.name, ''), COALESCE(c.description, '')
			FROM products p
			LEFT JOIN categories c ON p.category_id = c.id
			ORDER BY p.id`

	rows, err := repo.db.QueryContext(ctx, query)

	if err != nil {
		log.Printf("GetAll Query Error: %v\n", err)
		return nil , err 
	}

	defer rows.Close()

	products := make([]dto.ProductResponse, 0)

	for rows.Next(){
		var p dto.ProductResponse
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName, &p.CategoryDescription)
		if err != nil {
			log.Printf("GetAll Scan Error: %v\n", err)
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		log.Printf("GetAll rows.Err: %v\n", err)
		return nil, err
	}

	return products, nil
}

func (repo *ProductRepository) Create(ctx context.Context, product *models.Product) error{
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRowContext(ctx, query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)

	return err
}

func (repo *ProductRepository) Update(ctx context.Context, product *models.Product) error{
	query := "UPDATE products SET name=$1, price=$2, stock=$3, category_id=$4 WHERE id=$5"
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
	query := "DELETE FROM products WHERE id = $1"
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
	query := `SELECT p.id, p.name, p.price, p.stock, p.category_id,
			COALESCE(c.name, ''), COALESCE(c.description, '')
			FROM products p
			LEFT JOIN categories c ON p.category_id = c.id
			WHERE p.id = $1`

	var p dto.ProductResponse
	err := repo.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName, &p.CategoryDescription)
	if err == sql.ErrNoRows {
		return nil, errors.New("Product not found")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}