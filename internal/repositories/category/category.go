package category

import (
	"context"
	"ne-project/internal/models"

	"github.com/jmoiron/sqlx"
)

type CategoryRepositoryItf interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	Create(ctx context.Context, category *models.Category) (error)
	GetByID(ctx context.Context, id int) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id int) error
}

type CategoryRepository struct {
	db  *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepositoryItf {
	return &CategoryRepository{
		db: db,
	}
}