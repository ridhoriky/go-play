package product

import (
	"encoding/json"
	"log"
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *productHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.GetProductsQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrInvalidReqBody,
		})
		return
	}
	products, err := h.productService.GetAllProducts(ctx, &req)
	if err != nil {
		helpers.ResponseError(c.Writer, err)
		return
	}
	helpers.ResponseSuccess(c.Writer, http.StatusOK, "Success", products)
}

func (h *productHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var p dto.CreateProductRequest

	if err := c.ShouldBindJSON(&p); err != nil {
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrInvalidReqBody,
		})
		return
	}

	if _, err := h.productService.CreateProduct(ctx, &p); err != nil {
		helpers.ResponseError(c.Writer, err)
		return
	}
	helpers.ResponseSuccess(c.Writer, http.StatusCreated, "Product created successfully", p)
}

func (h *productHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrInvalidProductID,
		})
		return
	}
	product, err := h.productService.GetProductByID(ctx, id.String())
	if err != nil {
		helpers.ResponseError(c.Writer, err)
		return
	}
	helpers.ResponseSuccess(c.Writer, http.StatusOK, "Success", product)
}

func (h *productHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrInvalidProductID,
		})
		return
	}
	var p dto.UpdateProductRequest

	if err := c.ShouldBindJSON(&p); err != nil {
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrInvalidReqBody,
		})
		return
	}

	product, err := h.productService.UpdateProduct(ctx, id.String(), &p)
	log.Println(product, err)
	if err != nil {
		helpers.ResponseError(c.Writer, err)
		return
	}
	helpers.ResponseSuccess(c.Writer, http.StatusOK, "Product updated successfully", product)
}

func (h *productHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrInvalidProductID,
		})
		return
	}
	if err := h.productService.DeleteProduct(ctx, id.String()); err != nil {
		helpers.ResponseError(c.Writer, err)
		return
	}
	helpers.ResponseSuccess(c.Writer, http.StatusOK, "Product deleted successfully", nil)
}

func (h *productHandler) CreateMultiple(c *gin.Context) {
	ctx := c.Request.Context()
	var products []entity.Product
	if err := json.NewDecoder(c.Request.Body).Decode(&products); err != nil {
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrInvalidReqBody,
		})
		return
	}
	if len(products) == 0 {
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrNoProductCreated,
		})
		return
	}
	responses, err := h.productService.CreateMultipleProducts(ctx, products)
	if err != nil {
		helpers.ResponseError(c.Writer, err)
		return
	}
	helpers.ResponseSuccess(c.Writer, http.StatusCreated, "Products created successfully", responses)
}
