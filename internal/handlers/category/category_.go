package category

import (
	"context"
	"encoding/json"
	"ne-project/internal/dto"
	"net/http"
	"strconv"
)

// HandleCategories - GET /api/categories
func (h *categoryHandler) HandleCategories( w http.ResponseWriter, r *http.Request) {
	
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		h.GetAll(ctx, w, r)
	case http.MethodPost:
		h.Create(ctx, w, r)
	default:
		dto.ResponseError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *categoryHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryService.GetAllCategories(ctx) 
	if err != nil {
		dto.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	dto.ResponseSuccess(w, http.StatusOK, "Success", categories)
}

func (h *categoryHandler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var c dto.CategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		dto.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := c.Validate(); err != nil {
		dto.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.categoryService.CreateCategory(ctx, &c); err != nil {
		dto.ResponseError(w, http.StatusInternalServerError, "Failed to create category")
		return
	}
	dto.ResponseSuccess(w, http.StatusCreated, "Category created successfully", c)
}

// HandleCategoryByID - GET, PUT, DELETE /api/categories/{id}
func (h *categoryHandler) HandleCategoryByID( w http.ResponseWriter, r *http.Request) {
	
	ctx := r.Context()
	switch r.Method {
	case http.MethodGet:
		h.GetByID(ctx, w, r)
	case http.MethodPut:
		h.Update(ctx, w, r)
	case http.MethodDelete:
		h.Delete(ctx, w, r)
	default:
		dto.ResponseError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}

}

func (h *categoryHandler) GetByID(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := r.URL.Path[len("/api/categories/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.ResponseError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}
	category, err := h.categoryService.GetCategoryByID(ctx, id)
	if err != nil {
		dto.ResponseError(w, http.StatusInternalServerError, "Failed to fetch category")
		return
	}
	dto.ResponseSuccess(w, http.StatusOK, "Success", category)
}

func (h *categoryHandler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := r.URL.Path[len("/api/categories/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.ResponseError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}
	var c dto.CategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		dto.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := c.Validate(); err != nil {
		dto.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	c.ID = id
	if err := h.categoryService.UpdateCategory(ctx, &c); err != nil {
		dto.ResponseError(w, http.StatusInternalServerError, "Failed to update category")
		return
	}
	dto.ResponseSuccess(w, http.StatusOK, "Category updated successfully", c)
}

func (h *categoryHandler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := r.URL.Path[len("/api/categories/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.ResponseError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}
	if err := h.categoryService.DeleteCategory(ctx, id); err != nil {
		dto.ResponseError(w, http.StatusInternalServerError, "Failed to delete category")
		return
	}
	dto.ResponseSuccess(w, http.StatusOK, "Category deleted successfully", nil)
}

