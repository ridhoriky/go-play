package rest

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/product"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductHandler struct {
	productService product.ProductServiceItf
}

func NewProductHandler(productService product.ProductServiceItf) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// GetProducts godoc
// @Summary      List Product
// @Description  Get list of products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        search       		query   	string  false  "Filter by product name"
// @Param        category_id  		query   	string  false  "Filter by category ID (UUID)"
// @Param        min_price    		query   	number  false  "Minimum product price"
// @Param        max_price    		query   	number  false  "Maximum product price"
// @Param        in_stock     		query   	bool    false  "Filter products that have stock"
// @Param		 page				query		int		false	"Page number"	default(1)
// @Param		 limit				query		int		false	"Page size"		default(10)
// @Param		 sort_by			query		string	false	"Sort by field"
// @Param		 sort_dir			query		string	false	"Sort direction (asc/desc)"	default(asc)
// @Success      200  {object}  	dto.APIResponse{data=[]entity.ProductWithCategory}
// @Failure      400  {object} 		dto.APIResponse
// @Failure      404  {object}  	dto.APIResponse
// @Router       /products [get]
func (h *ProductHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.GetProductsQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidQueryParams))
		return
	}
	products, err := h.productService.GetAllProducts(ctx, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", products)
}

// CreateProduct godoc
// @Summary      Create Single product
// @Description  Create a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param		 product	body		dto.CreateProductRequest		true	"Product data"
// @Success      201		{object}  	dto.APIResponse{data=entity.Product}
// @Failure      400  		{object} 	dto.APIResponse
// @Failure      404  		{object}  	dto.APIResponse
// @Router       /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	var p dto.CreateProductRequest

	if err := c.ShouldBindJSON(&p); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	createdProduct, err := h.productService.CreateProduct(ctx, &p)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Product created successfully", createdProduct)
}

// GetProductByID godoc
// @Summary      Get Single product
// @Description  Get data product by id
// @Tags         products
// @Produce      json
// @Param 		 id   path 		string true "Product ID (UUID)" format(uuid)
// @Success      200  {object}  dto.APIResponse{data=entity.Product}
// @Failure      400  {object} 	dto.APIResponse
// @Failure      404  {object}  dto.APIResponse
// @Router       /products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidProductID))
		return
	}
	product, err := h.productService.GetProductByID(ctx, id.String())
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", product)
}

// UpdateProduct godoc
// @Summary      Update Product
// @Description  Update a product by their id
// @Tags         products
// @Accept       json
// @Produce      json
// @Param 		 id 	  path 	string true "Product ID (UUID)" format(uuid)
// @Param 		 product  body 	dto.UpdateProductRequest true "Product data"
// @Success      200  {object}  dto.APIResponse{data=entity.Product}
// @Failure      400  {object} 	dto.APIResponse
// @Failure      404  {object}  dto.APIResponse
// @Router       /products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidProductID))
		return
	}
	var p dto.UpdateProductRequest

	if err := c.ShouldBindJSON(&p); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	updatedProduct, err := h.productService.UpdateProduct(ctx, id.String(), &p)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Product updated successfully", updatedProduct)
}

// DeleteProduct godoc
// @Summary Delete Product
// @Description  Delete a product by their id
// @Tags         products
// @Accept       json
// @Produce      json
// @Param 		 id   path 	string true "Product ID (UUID)" format(uuid)
// @Success      200  {object}  dto.APIResponse
// @Failure      400  {object} 	dto.APIResponse
// @Failure      404  {object}  dto.APIResponse
// @Router       /products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidProductID))
		return
	}
	if err := h.productService.DeleteProduct(ctx, id.String()); err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Product deleted successfully", nil)
}

// CreateProduct godoc
// @Summary      Create Multiple product
// @Description  Create multiple products in a single request
// @Tags         products
// @Accept       json
// @Produce      json
// @Param 		 product 	body 		dto.CreateMultipleProducts true "Product data"
// @Success 	 201 		{object} 	dto.APIResponse{data=[]entity.Product}
// @Failure      400  		{object} 	dto.APIResponse
// @Failure      404  		{object}  	dto.APIResponse
// @Router       /products/bulk [post]
func (h *ProductHandler) CreateMultiple(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.CreateMultipleProducts
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	var products []entity.Product
	for _, p := range req.Data {
		products = append(products, entity.Product{
			Name:       p.Name,
			Price:      decimal.NewFromFloat(p.Price),
			CategoryID: p.CategoryID,
		})
	}

	responses, err := h.productService.CreateMultipleProducts(ctx, products)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Products created successfully", responses)
}
