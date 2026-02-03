package category

import (
	"context"
	"ne-project/internal/dto"
	"ne-project/internal/models"
)

func (s *categoryService) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	return s.categoryRepository.GetAll(ctx)
}
func (s *categoryService) CreateCategory(ctx context.Context, data *dto.CategoryDTO) error {
	return s.categoryRepository.Create(ctx, data.ToModelPtr())
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id int) (*models.Category, error) {
	return s.categoryRepository.GetByID(ctx, id)
}

func (s *categoryService) UpdateCategory(ctx context.Context, category *dto.CategoryDTO) error {
	return s.categoryRepository.Update(ctx, category.ToModelPtr())
}

func (s *categoryService) DeleteCategory(ctx context.Context, id int) error {
	return s.categoryRepository.Delete(ctx, id)
}