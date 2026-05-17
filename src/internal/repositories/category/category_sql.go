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

	"github.com/rs/zerolog"
)

func (repo *categoryRepository) GetAll(ctx context.Context, filter *dto.GetCategoriesQuery) ([]entity.Category, int, error) {
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

	args = append(args, filter.Limit, offset)

	rows, err := repo.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err find category with query")
		return nil, 0, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			zerolog.Ctx(ctx).Error().Err(closeErr).Msg("failed to close rows")
		}
	}()

	categories := []entity.Category{}
	var total int

	for rows.Next() {

		var c entity.Category
		var rowCount int

		err = rows.Scan(
			&c.ID,
			&c.Name,
			&c.Description,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.DeletedAt,
			&rowCount,
		)

		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Str("categoryID", c.ID).Msg("err mapping category row")
			return nil, 0, err
		}

		categories = append(categories, c)

		// assign value total from firt row
		if total == 0 {
			total = rowCount
		}
	}

	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err iterating category rows")
		return nil, 0, err
	}

	return categories, total, nil
}

func buildCategoryFilters(filter *dto.GetCategoriesQuery) (string, []any) {

	query := ""
	args := []any{}

	if filter.Search != "" {
		query += " AND c.name ILIKE $1"
		args = append(args, "%"+filter.Search+"%")
	}

	if !filter.IncludeDeleted {
		query += " AND c.deleted_at IS NULL"
	}

	return query, args
}

func (repo *categoryRepository) Create(ctx context.Context, category *entity.Category) error {
	err := repo.db.QueryRowContext(ctx, createCategoryQuery, category.Name, category.Description).Scan(&category.ID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err to create category")
		return err
	}
	return nil
}

func (repo *categoryRepository) GetByID(ctx context.Context, id string) (*entity.Category, error) {
	var c entity.Category
	err := repo.db.QueryRowContext(ctx, getCategoryByIDQuery, id).Scan(&c.ID, &c.Name, &c.Description)
	if errors.Is(err, sql.ErrNoRows) {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err category id not found")
		return nil, dto.NewError(http.StatusNotFound, preference.ErrCategoryNotFound)
	}

	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err get single category")
		return nil, err
	}

	return &c, nil
}

func (repo *categoryRepository) Update(ctx context.Context, id string, category *entity.Category) error {
	result, err := repo.db.ExecContext(ctx, updateCategoryQuery, category.Name, category.Description, id)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err update category")
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err get rows affected")
		return err
	}
	if rows == 0 {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err not found update category")
		return dto.NewError(http.StatusNotFound, preference.ErrCategoryNotFound)
	}
	return nil
}

func (repo *categoryRepository) Delete(ctx context.Context, id string) error {
	result, err := repo.db.ExecContext(ctx, deleteCategoryQuery, id)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err delete category")
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err get rows affected")
		return err
	}

	if rows == 0 {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err not found delete category")
		return dto.NewError(http.StatusNotFound, preference.ErrCategoryNotFound)
	}

	return err

}
