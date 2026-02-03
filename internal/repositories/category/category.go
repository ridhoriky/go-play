package category

import (
	"context"
	"database/sql"
	"ne-project/internal/models"
)

type CategoryRepositoryItf interface {
	GetAll(ctx context.Context) ([]models.Category, error)
	Create(ctx context.Context, category *models.Category) (error)
	GetByID(ctx context.Context, id int) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id int) error
}

type CategoryRepository struct {
	db  *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepositoryItf {
	return &CategoryRepository{
		db: db,
	}
}