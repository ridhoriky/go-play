package product

import (
	"encoding/json"
	"ne-project/internal/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *productHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	products, err := h.productService.GetAllProducts(ctx) 
	if err != nil {
		dto.ResponseError(c.Writer, http.StatusInternalServerError, err.Error())
		return
	}
	dto.ResponseSuccess(c.Writer, http.StatusOK, "Success", products)
}

func (h *productHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var p dto.ProductDTO
	if err := json.NewDecoder(c.Request.Body).Decode(&p); err != nil {
		dto.ResponseError(c.Writer, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := p.Validate(); err != nil {
		dto.ResponseError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.productService.CreateProduct(ctx, &p); err != nil {
		dto.ResponseError(c.Writer, http.StatusInternalServerError, "Failed to create product")
		return
	}
	dto.ResponseSuccess(c.Writer, http.StatusCreated, "Product created successfully", p)
}

func (h *productHandler) GetByID(c *gin.Context) {
	ctx:= c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.ResponseError(c.Writer, http.StatusBadRequest, "Invalid product ID")
		return
	}
	product, err := h.productService.GetProductByID(ctx, id)
	if err != nil {
		dto.ResponseError(c.Writer, http.StatusInternalServerError, "Failed to fetch product")
		return
	}
	dto.ResponseSuccess(c.Writer, http.StatusOK, "Success", product)
}

func (h *productHandler) Update(c *gin.Context) {
	ctx:= c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.ResponseError(c.Writer, http.StatusBadRequest, "Invalid product ID")
		return
	}
	var p dto.ProductDTO
	if err := json.NewDecoder(c.Request.Body).Decode(&p); err != nil {
		dto.ResponseError(c.Writer, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := p.Validate(); err != nil {
		dto.ResponseError(c.Writer, http.StatusBadRequest, err.Error())
		return
	}
	p.ID = id
	if err := h.productService.UpdateProduct(ctx, &p); err != nil {
		dto.ResponseError(c.Writer, http.StatusInternalServerError, "Failed to update product")
		return
	}
	dto.ResponseSuccess(c.Writer, http.StatusOK, "Product updated successfully", p)
}

func (h *productHandler) Delete(c *gin.Context) {
	ctx:= c.Request.Context()
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		dto.ResponseError(c.Writer, http.StatusBadRequest, "Invalid product ID")
		return
	}
	if err := h.productService.DeleteProduct(ctx, id); err != nil {
		dto.ResponseError(c.Writer, http.StatusInternalServerError, "Failed to delete product")
		return
	}
	dto.ResponseSuccess(c.Writer, http.StatusOK, "Product deleted successfully", nil)
}