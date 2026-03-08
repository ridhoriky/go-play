package category

import (
	"context"
	"errors"

	"ne-project/src/internal/models/entity"
)

func (repo *CategoryRepository) GetAll(ctx context.Context) ([]entity.Category, error) {
	rows, err := repo.db.QueryContext(ctx, getAllCategoriesQuery)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]entity.Category, 0)
	for rows.Next() {
		var c entity.Category
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

func (repo *CategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	err := repo.db.QueryRowContext(ctx, createCategoryQuery, category.Name, category.Description).Scan(&category.ID)
	return err
}

func (repo *CategoryRepository) GetByID(ctx context.Context, id string) (*entity.Category, error) {
	var c entity.Category
	err := repo.db.QueryRowContext(ctx, getCategoryByIDQuery, id).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (repo *CategoryRepository) Update(ctx context.Context, category *entity.Category) error {
	result, err := repo.db.ExecContext(ctx, updateCategoryQuery, category.Name, category.Description, category.ID)
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

func (repo *CategoryRepository) Delete(ctx context.Context, id string) error {
	result, err := repo.db.ExecContext(ctx, deleteCategoryQuery, id)
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
