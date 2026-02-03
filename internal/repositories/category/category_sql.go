package category

import (
	"context"
	"errors"
	"ne-project/internal/models"
)

func (repo *CategoryRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	query := "SELECT id, name, description FROM categories"
	rows, err := repo.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name, &c.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (repo *CategoryRepository) Create(ctx context.Context, category *models.Category) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRowContext(ctx, query, category.Name, category.Description).Scan(&category.ID)
	return err
}

func (repo *CategoryRepository) GetByID(ctx context.Context, id int) (*models.Category, error) {
	query := "SELECT id, name, description FROM categories WHERE id=$1"
	var c models.Category
	err := repo.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (repo *CategoryRepository) Update(ctx context.Context, category *models.Category) error {
	query := "UPDATE categories SET name=$1, description=$2 WHERE id=$3"
	result, err := repo.db.ExecContext(ctx, query, category.Name, category.Description, category.ID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Category not found")
	}
	return nil
}

func (repo *CategoryRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM categories WHERE id=$1"
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