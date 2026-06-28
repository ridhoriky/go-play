package category

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/repositories/category"
)

type CategoryServiceItf interface {
	GetCategoryTree(ctx context.Context) ([]dto.CategoryTreeNode, error)
	GetCategoryByID(ctx context.Context, id string) (*dto.CategoryDetailResponse, error)
	CreateCategory(ctx context.Context, category *dto.CreateCategoryRequest) (*entity.Category, error)
	UpdateCategory(ctx context.Context, id string, category *dto.UpdateCategoryRequest) (*entity.Category, error)
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
