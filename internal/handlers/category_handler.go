package handlers

import (
	"encoding/json"
	"ne-project/internal/dto"
	"ne-project/internal/services"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(service *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// HandleCategories - GET /api/categories
func (h *CategoryHandler) HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAll()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondSuccess(w, http.StatusOK, "Success", categories)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c dto.CategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := c.Validate(); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	catModel := c.ToModel()
	if err := h.service.Create(&catModel); err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to create category")
		return
	}
	RespondSuccess(w, http.StatusCreated, "Category created successfully", catModel)
}

// HandleCategoryByID - GET, PUT, DELETE /api/categories/{id}
func (h *CategoryHandler) HandleCategoryByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetByID(w, r)
	case http.MethodPut:
		h.Update(w, r)
	case http.MethodDelete:
		h.Delete(w, r)
	default:
		RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}

}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := r.URL.Path[len("/api/categories/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}
	category, err := h.service.GetByID(id)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to fetch category")
		return
	}
	RespondSuccess(w, http.StatusOK, "Success", category)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := r.URL.Path[len("/api/categories/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}
	var c dto.CategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := c.Validate(); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	catModel := c.ToModel()
	catModel.ID = id
	if err := h.service.Update(&catModel); err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to update category")
		return
	}
	RespondSuccess(w, http.StatusOK, "Category updated successfully", catModel)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := r.URL.Path[len("/api/categories/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}
	if err := h.service.Delete(id); err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to delete category")
		return
	}
	RespondSuccess(w, http.StatusOK, "Category deleted successfully", nil)
}