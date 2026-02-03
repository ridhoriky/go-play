package category

import (
	"context"
	"ne-project/internal/dto"
	"ne-project/internal/models"
	"ne-project/internal/repositories/category"
)

type CategoryServiceItf interface {
	GetAllCategories(ctx context.Context) ([]models.Category, error)
	GetCategoryByID(ctx context.Context, id int) (*models.Category, error)
	CreateCategory(ctx context.Context, category *dto.CategoryDTO) error
	UpdateCategory(ctx context.Context, category *dto.CategoryDTO) error
	DeleteCategory(ctx context.Context, id int) error
}


type categoryService struct {
	categoryRepository category.CategoryRepositoryItf
}

func NewCategoryService(categoryRepository category.CategoryRepositoryItf) CategoryServiceItf {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}

