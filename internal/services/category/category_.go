package category

import (
	"context"

	"ne-project/internal/models/dto"
	"ne-project/internal/models/entity"
)

func (s *categoryService) GetAllCategories(ctx context.Context) ([]entity.Category, error) {
	return s.categoryRepository.GetAll(ctx)
}
func (s *categoryService) CreateCategory(ctx context.Context, data *dto.CategoryDTO) error {
	return s.categoryRepository.Create(ctx, data.ToModelPtr())
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id string) (*entity.Category, error) {
	return s.categoryRepository.GetByID(ctx, id)
}

func (s *categoryService) UpdateCategory(ctx context.Context, category *dto.CategoryDTO) error {
	return s.categoryRepository.Update(ctx, category.ToModelPtr())
}

func (s *categoryService) DeleteCategory(ctx context.Context, id string) error {
	return s.categoryRepository.Delete(ctx, id)
}
