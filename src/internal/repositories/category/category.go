package category

import (
	"context"

	"ne-project/src/internal/models/entity"

	"github.com/jmoiron/sqlx"
)

type CategoryRepositoryItf interface {
	HasProducts(ctx context.Context, categoryID string) (bool, error)
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id string) (*entity.Category, error)
	Update(ctx context.Context, id string, category *entity.Category) error
	Delete(ctx context.Context, id string) error
	GetCategoryTree(ctx context.Context) ([]entity.CategoryWithCount, error)
	GetByParentID(ctx context.Context, parentID string) ([]entity.Category, error)
}

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepositoryItf {
	return &categoryRepository{
		db: db,
	}
}
