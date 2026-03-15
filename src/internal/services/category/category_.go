package category

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"net/http"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"
)

func (s *categoryService) GetAllCategories(ctx context.Context, req *dto.GetCategoriesQuery) (*dto.CategoryListResponse, error) {
	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	categories, total, err := s.categoryRepository.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	res := make([]dto.CategoryResponse, 0, len(categories))

	for _, c := range categories {
		resCategory := dto.CategoryResponse{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
			DeletedAt:   c.DeletedAt,
		}

		res = append(res, resCategory)
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	resp := &dto.CategoryListResponse{
		Data: res,
		Meta: dto.PaginationMeta{
			Total:      total,
			Page:       req.Page,
			Limit:      req.Limit,
			TotalPages: totalPages,
		},
	}

	return resp, nil
}
func (s *categoryService) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*entity.Category, error) {

	c := &entity.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.categoryRepository.Create(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id string) (*entity.Category, error) {
	return s.categoryRepository.GetByID(ctx, id)
}

func (s *categoryService) UpdateCategory(ctx context.Context, id string, req *dto.UpdateCategoryRequest) error {
	existingCategory, err := s.categoryRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &dto.Error{
				Code:    http.StatusBadRequest,
				Message: preference.ErrInvalidCategoryID,
			}
		}
		return err
	}

	existingCategory.Name = req.Name

	return s.categoryRepository.Update(ctx, id, existingCategory)
}

func (s *categoryService) DeleteCategory(ctx context.Context, id string) error {
	return s.categoryRepository.Delete(ctx, id)
}
