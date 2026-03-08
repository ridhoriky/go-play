package category

import (
	"context"

	"ne-project/src/internal/models/entity"

	"github.com/jmoiron/sqlx"
)

type CategoryRepositoryItf interface {
	GetAll(ctx context.Context) ([]entity.Category, error)
	Create(ctx context.Context, category *entity.Category) error
	GetByID(ctx context.Context, id string) (*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id string) error
}

type CategoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepositoryItf {
	return &CategoryRepository{
		db: db,
	}
}
