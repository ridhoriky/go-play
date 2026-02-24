package category

import (
	"context"

	"ne-project/internal/models/dto"
	"ne-project/internal/models/entity"
	"ne-project/internal/repositories/category"
)

type CategoryServiceItf interface {
	GetAllCategories(ctx context.Context) ([]entity.Category, error)
	GetCategoryByID(ctx context.Context, id string) (*entity.Category, error)
	CreateCategory(ctx context.Context, category *dto.CategoryDTO) error
	UpdateCategory(ctx context.Context, category *dto.CategoryDTO) error
	DeleteCategory(ctx context.Context, id string) error
}

type categoryService struct {
	categoryRepository category.CategoryRepositoryItf
}

func NewCategoryService(categoryRepository category.CategoryRepositoryItf) CategoryServiceItf {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}
