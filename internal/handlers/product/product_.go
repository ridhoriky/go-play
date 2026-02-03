package product

import (
	"context"
	"encoding/json"
	"ne-project/internal/dto"
	"net/http"
	"strconv"
)

// HandleGetAllProducts - GET /api/products
func (h *productHandler) HandleProducts( w http.ResponseWriter, r *http.Request) {
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

func (h *productHandler) GetAll(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.GetAllProducts(ctx) 
	if err != nil {
		dto.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}
	dto.ResponseSuccess(w, http.StatusOK, "Success", products)
}

func (h *productHandler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var c dto.ProductDTO
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		dto.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := c.Validate(); err != nil {
		dto.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.productService.CreateProduct(ctx, &c); err != nil {
		dto.ResponseError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}
	dto.ResponseSuccess(w, http.StatusCreated, "Product created successfully", c)
}

// HandleProductByID - GET, PUT, DELETE /api/products/{id}
func (h *productHandler) HandleProductByID( w http.ResponseWriter, r *http.Request) {
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

func (h *productHandler) GetByID(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/products/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.ResponseError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	product, err := h.productService.GetProductByID(ctx, id)
	if err != nil {
		dto.ResponseError(w, http.StatusInternalServerError, "Failed to fetch product")
		return
	}
	dto.ResponseSuccess(w, http.StatusOK, "Success", product)
}

func (h *productHandler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/products/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.ResponseError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	var p dto.ProductDTO
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		dto.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := p.Validate(); err != nil {
		dto.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	p.ID = id
	if err := h.productService.UpdateProduct(ctx, &p); err != nil {
		dto.ResponseError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}
	dto.ResponseSuccess(w, http.StatusOK, "Product updated successfully", p)
}

func (h *productHandler) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/products/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.ResponseError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}
	if err := h.productService.DeleteProduct(ctx, id); err != nil {
		dto.ResponseError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}
	dto.ResponseSuccess(w, http.StatusOK, "Product deleted successfully", nil)
}