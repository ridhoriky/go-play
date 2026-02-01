package handlers

import (
	"encoding/json"
	"ne-project/internal/dto"
	"ne-project/internal/services"
	"net/http"
	"strconv"
	"strings"
)
type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// HandleProducts - GET /api/products
func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetAll(w, r)
	case http.MethodPost:
		h.Create(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to fetch product")
		return
	}

	RespondSuccess(w, http.StatusOK, "Success", products)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p dto.ProductDTO
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := p.Validate(); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	prodModel := p.ToModel()
	if err := h.service.Create(&prodModel); err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	RespondSuccess(w, http.StatusCreated, "Product created successfully", prodModel)
}

// HandleProductByID - GET/PUT/DELETE /api/product/{id}
func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
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

// GetByID - GET /api/product/{id}
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to fetch product")
		return
	}

	RespondSuccess(w, http.StatusOK, "Success", product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	var p dto.ProductDTO
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := p.Validate(); err != nil {
		RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	prodModel := p.ToModel()
	prodModel.ID = id
	if err := h.service.Update(&prodModel); err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	RespondSuccess(w, http.StatusOK, "Product updated successfully", prodModel)
}

// Delete - DELETE /api/product/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	RespondSuccess(w, http.StatusOK, "Product deleted successfully", nil)
}