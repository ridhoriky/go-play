package category

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"
)

func (repo *CategoryRepository) GetAll(ctx context.Context, filter *dto.GetCategoriesQuery) ([]entity.Category, int, error) {
	filterQuery, args := buildCategoryFilters(filter)
	dataQuery := getAllCategoriesQuery + filterQuery

	sortBy := "c.created_at"
	allowedSortColumns := map[string]string{
		"name":       "c.name",
		"created_at": "c.created_at",
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

	categories := []entity.Category{}

	for rows.Next() {

		var c entity.Category

		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.DeletedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		categories = append(categories, c)
	}

	countQuery := countCategoriesQuery + filterQuery

	var total int

	err = repo.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func buildCategoryFilters(filter *dto.GetCategoriesQuery) (string, []interface{}) {

	query := ""
	args := []interface{}{}
	argPos := 1

	// search
	if filter.Search != "" {
		query += fmt.Sprintf(" AND c.name ILIKE $%d", argPos)
		args = append(args, "%"+filter.Search+"%")
		argPos++
	}

	// isDeleted filter
	if filter.IncludeDeleted {
		query += " AND c.deleted_at IS NOT NULL"
	} else {
		query += " AND c.deleted_at IS NULL"
	}

	return query, args
}

func (repo *CategoryRepository) Create(ctx context.Context, category *entity.Category) error {
	err := repo.db.QueryRowContext(ctx, createCategoryQuery, category.Name, category.Description).Scan(&category.ID)
	return err
}

func (repo *CategoryRepository) GetByID(ctx context.Context, id string) (*entity.Category, error) {
	var c entity.Category
	err := repo.db.QueryRowContext(ctx, getCategoryByIDQuery, id).Scan(&c.ID, &c.Name, &c.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, &dto.Error{
			Code:    http.StatusNotFound,
			Message: preference.ErrCategoryNotFound,
		}
	}

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (repo *CategoryRepository) Update(ctx context.Context, id string, category *entity.Category) error {
	result, err := repo.db.ExecContext(ctx, updateCategoryQuery, category.Name, category.Description, id)
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
			Message: preference.ErrCategoryNotFound,
		}
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
		return &dto.Error{
			Code:    http.StatusNotFound,
			Message: preference.ErrCategoryNotFound,
		}
	}

	return err

}
