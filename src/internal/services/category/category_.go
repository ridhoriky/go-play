package category

import (
	"context"
	"net/http"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"

	"github.com/rs/zerolog"
)

func (s *categoryService) GetCategoryTree(ctx context.Context) ([]dto.CategoryTreeNode, error) {
	categories, err := s.categoryRepository.GetCategoryTree(ctx)
	if err != nil {
		return nil, err
	}
	return buildCategoryTree(categories, nil), nil
}

func buildCategoryTree(categories []entity.CategoryWithCount, parentID *string) []dto.CategoryTreeNode {
	var tree []dto.CategoryTreeNode
	for i := range categories {
		c := categories[i]
		if (parentID == nil && c.ParentID == nil) || (parentID != nil && c.ParentID != nil && *parentID == *c.ParentID) {
			node := dto.CategoryTreeNode{
				ID:           c.ID,
				Name:         c.Name,
				Description:  c.Description,
				ImageURL:     c.ImageURL,
				SortOrder:    c.SortOrder,
				ProductCount: c.ProductCount,
				Children:     buildCategoryTree(categories, &c.ID),
			}
			tree = append(tree, node)
		}
	}
	return tree
}

func (s *categoryService) CreateCategory(ctx context.Context, req *dto.CreateCategoryRequest) (*entity.Category, error) {
	if req.Name == "" {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrCategoryNameRequired)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrCategoryNameRequired)
	}

	c := &entity.Category{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		ImageURL:    req.ImageURL,
		SortOrder:   req.SortOrder,
	}

	if err := s.categoryRepository.Create(ctx, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *categoryService) GetCategoryByID(ctx context.Context, id string) (*dto.CategoryDetailResponse, error) {
	category, err := s.categoryRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Fetch children
	childrenEntity, err := s.categoryRepository.GetByParentID(ctx, category.ID)
	if err != nil {
		return nil, err
	}
	var children []dto.CategoryResponse
	for _, c := range childrenEntity {
		children = append(children, dto.CategoryResponse{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			ParentID:    c.ParentID,
			ImageURL:    c.ImageURL,
			SortOrder:   c.SortOrder,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		})
	}

	// Breadcrumb (go up the tree)
	var breadcrumb []dto.CategoryResponse
	currParent := category.ParentID
	for currParent != nil {
		parent, err := s.categoryRepository.GetByID(ctx, *currParent)
		if err != nil {
			break
		}
		breadcrumb = append([]dto.CategoryResponse{{
			ID:          parent.ID,
			Name:        parent.Name,
			Description: parent.Description,
			ParentID:    parent.ParentID,
			ImageURL:    parent.ImageURL,
			SortOrder:   parent.SortOrder,
			CreatedAt:   parent.CreatedAt,
			UpdatedAt:   parent.UpdatedAt,
		}}, breadcrumb...)
		currParent = parent.ParentID
	}

	return &dto.CategoryDetailResponse{
		CategoryResponse: dto.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			ParentID:    category.ParentID,
			ImageURL:    category.ImageURL,
			SortOrder:   category.SortOrder,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		},
		ProductsCount: 0, // Product count logic can be added later if needed
		Children:      children,
		Breadcrumb:    breadcrumb,
	}, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, id string, req *dto.UpdateCategoryRequest) (*entity.Category, error) {
	existingCategory, err := s.categoryRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name == "" {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrCategoryNameRequired)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrCategoryNameRequired)
	}

	existingCategory.Name = req.Name
	existingCategory.Description = req.Description
	existingCategory.ParentID = req.ParentID
	existingCategory.ImageURL = req.ImageURL
	existingCategory.SortOrder = req.SortOrder

	if err := s.categoryRepository.Update(ctx, id, existingCategory); err != nil {
		return nil, err
	}
	return existingCategory, nil
}

func (s *categoryService) DeleteCategory(ctx context.Context, id string) error {
	children, err := s.categoryRepository.GetByParentID(ctx, id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return dto.NewError(http.StatusBadRequest, "Cannot delete category that has sub-categories")
	}

	hasProducts, err := s.categoryRepository.HasProducts(ctx, id)
	if err != nil {
		return err
	}
	if hasProducts {
		return dto.NewError(http.StatusBadRequest, "Cannot delete category that has products")
	}

	return s.categoryRepository.Delete(ctx, id)
}
