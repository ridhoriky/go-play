package category

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"

	"github.com/rs/zerolog"
)

func (repo *categoryRepository) HasProducts(ctx context.Context, categoryID string) (bool, error) {
	var hasProducts bool
	err := repo.db.QueryRowContext(ctx, hasProductsQuery, categoryID).Scan(&hasProducts)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("category_id", categoryID).Msg("err check if category has products")
		return false, err
	}
	return hasProducts, nil
}

func (repo *categoryRepository) Create(ctx context.Context, category *entity.Category) error {
	err := repo.db.QueryRowContext(ctx, createCategoryQuery, category.Name, category.Description, category.ParentID, category.ImageURL, category.SortOrder).Scan(&category.ID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err to create category")
		return err
	}
	return nil
}

func (repo *categoryRepository) GetByID(ctx context.Context, id string) (*entity.Category, error) {
	var c entity.Category
	err := repo.db.QueryRowContext(ctx, getCategoryByIDQuery, id).Scan(&c.ID, &c.Name, &c.Description, &c.ParentID, &c.ImageURL, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)
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
	result, err := repo.db.ExecContext(ctx, updateCategoryQuery, category.Name, category.Description, category.ParentID, category.ImageURL, category.SortOrder, id)
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

func (repo *categoryRepository) GetCategoryTree(ctx context.Context) ([]entity.CategoryWithCount, error) {
	rows, err := repo.db.QueryContext(ctx, getCategoryTreeQuery)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err get category tree")
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			zerolog.Ctx(ctx).Error().Err(closeErr).Msg("failed to close rows")
		}
	}()

	var categories []entity.CategoryWithCount
	for rows.Next() {
		var c entity.CategoryWithCount
		err = rows.Scan(
			&c.ID, &c.Name, &c.Description, &c.ParentID, &c.ImageURL, &c.SortOrder,
			&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
			&c.ProductCount,
		)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("err scan category tree")
			return nil, err
		}
		categories = append(categories, c)
	}
	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err rows category tree")
		return nil, err
	}
	return categories, nil
}

func (repo *categoryRepository) GetByParentID(ctx context.Context, parentID string) ([]entity.Category, error) {
	rows, err := repo.db.QueryContext(ctx, getByParentIDQuery, parentID)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err get categories by parent")
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			zerolog.Ctx(ctx).Error().Err(closeErr).Msg("failed to close rows")
		}
	}()

	var categories []entity.Category
	for rows.Next() {
		var c entity.Category
		err = rows.Scan(
			&c.ID, &c.Name, &c.Description, &c.ParentID, &c.ImageURL, &c.SortOrder,
			&c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
		)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("err scan categories by parent")
			return nil, err
		}
		categories = append(categories, c)
	}
	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err rows categories by parent")
		return nil, err
	}
	return categories, nil
}
