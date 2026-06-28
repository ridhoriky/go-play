package rest

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/product"
	"ne-project/src/internal/services/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductHandler struct {
	productService product.ProductServiceItf
	storeService   store.StoreServiceItf
}

func NewProductHandler(productService product.ProductServiceItf, storeService store.StoreServiceItf) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		storeService:   storeService,
	}
}

// GetProducts godoc
// @Summary      List Product
// @Description  Get list of products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        q            		query   	string  false  "Filter by product name/description"
// @Param        category  			query   	string  false  "Filter by category ID (UUID)"
// @Param        store  			query   	string  false  "Filter by store ID (UUID)"
// @Param        min_price    		query   	number  false  "Minimum product price"
// @Param        max_price    		query   	number  false  "Maximum product price"
// @Param        rating             query       number  false  "Minimum product rating"
// @Param        in_stock     		query   	bool    false  "Filter products that have stock"
// @Param		 page				query		int		false	"Page number"	default(1)
// @Param		 limit				query		int		false	"Page size"		default(10)
// @Param		 sort				query		string	false	"Sort by field (newest, price_asc, price_desc, rating, popular)"
// @Success      200  {object}  	dto.APIResponse{data=dto.ProductListResponse}
// @Failure      400  {object} 		dto.APIResponse
// @Failure      404  {object}  	dto.APIResponse
// @Router       /products [get]
func (h *ProductHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetString("user_id")

	var req dto.GetProductsQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidQueryParams))
		return
	}
	products, err := h.productService.GetAllProducts(ctx, userID, &req)
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
	userID := c.GetString("user_id")

	s, err := h.storeService.GetMyStore(ctx, userID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	var p dto.CreateProductRequest

	if err = c.ShouldBindJSON(&p); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	createdProduct, err := h.productService.CreateProduct(ctx, s.ID, &p)
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
// @Success      200  {object}  dto.APIResponse{data=dto.ProductDetailResponse}
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
	userID := c.GetString("user_id")
	product, err := h.productService.GetProductByID(ctx, id.String(), userID)
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
	userID := c.GetString("user_id")

	s, err := h.storeService.GetMyStore(ctx, userID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidProductID))
		return
	}
	var p dto.UpdateProductRequest

	if err = c.ShouldBindJSON(&p); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	updatedProduct, err := h.productService.UpdateProduct(ctx, id.String(), s.ID, &p)
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
	userID := c.GetString("user_id")

	s, err := h.storeService.GetMyStore(ctx, userID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidProductID))
		return
	}
	if err := h.productService.DeleteProduct(ctx, id.String(), s.ID); err != nil {
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
	userID := c.GetString("user_id")

	s, err := h.storeService.GetMyStore(ctx, userID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	var req dto.CreateMultipleProducts
	if err = c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	var products []entity.Product
	for _, p := range req.Data {
		products = append(products, entity.Product{
			Name:        p.Name,
			Description: p.Description,
			Price:       decimal.NewFromFloat(p.Price),
			Stock:       p.Stock,
			CategoryID:  p.CategoryID,
		})
	}

	responses, err := h.productService.CreateMultipleProducts(ctx, s.ID, products)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Products created successfully", responses)
}

// GetMyProducts godoc
// @Summary      List My Product
// @Description  Get list of products belonging to the logged-in seller
// @Tags         seller-products
// @Produce      json
// @Success      200  {object}  dto.APIResponse{data=dto.ProductListResponse}
// @Router       /seller/products [get]
func (h *ProductHandler) GetMyProducts(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	s, err := h.storeService.GetMyStore(ctx, userID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	var req dto.GetProductsQuery
	if err = c.ShouldBindQuery(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidQueryParams))
		return
	}

	req.Store = s.ID
	products, err := h.productService.GetAllProducts(ctx, userID, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", products)
}

// GetStoreProducts godoc
// @Summary Get products of a store
// @Tags Stores
// @Success      200  {object}  dto.APIResponse{data=dto.ProductListResponse}
// @Param slug path string true "Store Slug"
// @Router /stores/{slug}/products [get]
func (h *ProductHandler) GetStoreProducts(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	s, err := h.storeService.GetStoreBySlug(ctx, slug)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	var req dto.GetProductsQuery
	if err = c.ShouldBindQuery(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidQueryParams))
		return
	}

	userID := c.GetString("user_id")
	req.Store = s.ID
	products, err := h.productService.GetAllProducts(ctx, userID, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", products)
}

// GetSellerProductDetail godoc
// @Summary      Get Seller Product Detail
// @Description  Get product detail for seller including recent orders
// @Tags         seller-products
// @Produce      json
// @Security     BearerAuth
// @Param 		 id   path 		string true "Product ID (UUID)" format(uuid)
// @Success      200  {object}  dto.APIResponse{data=dto.SellerProductDetailResponse}
// @Failure      400  {object} 	dto.APIResponse
// @Failure      404  {object}  dto.APIResponse
// @Router       /seller/products/{id} [get]
func (h *ProductHandler) GetSellerProductDetail(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	s, err := h.storeService.GetMyStore(ctx, userID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidProductID))
		return
	}

	product, err := h.productService.GetSellerProductDetail(ctx, id.String(), s.ID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", product)
}

// UpdateSellerProduct godoc
// @Summary      Update Seller Product
// @Description  Update a product by their id
// @Tags         seller-products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param 		 id 	  path 	string true "Product ID (UUID)" format(uuid)
// @Param 		 product  body 	dto.UpdateProductRequest true "Product data"
// @Success      200  {object}  dto.APIResponse{data=entity.Product}
// @Failure      400  {object} 	dto.APIResponse
// @Failure      404  {object}  dto.APIResponse
// @Router       /seller/products/{id} [put]
func (h *ProductHandler) UpdateSellerProduct(c *gin.Context) {
	// Re-using the update logic since it validates store ID
	h.Update(c)
}

// DeleteSellerProduct godoc
// @Summary Delete Seller Product
// @Description  Delete a product by their id with active order validation
// @Tags         seller-products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param 		 id   path 	string true "Product ID (UUID)" format(uuid)
// @Success      200  {object}  dto.APIResponse
// @Failure      400  {object} 	dto.APIResponse
// @Failure      404  {object}  dto.APIResponse
// @Router       /seller/products/{id} [delete]
func (h *ProductHandler) DeleteSellerProduct(c *gin.Context) {
	// Re-using the delete logic since it validates store ID
	h.Delete(c)
}

// AddProductImage godoc
// @Summary      Add Product Image
// @Description  Add an image to a product
// @Tags         seller-products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param 		 id   path 	string true "Product ID (UUID)" format(uuid)
// @Param 		 image body dto.AddProductImageRequest true "Image data"
// @Success      201  {object}  dto.APIResponse{data=entity.ProductImage}
// @Failure      400  {object} 	dto.APIResponse
// @Router       /seller/products/{id}/images [post]
func (h *ProductHandler) AddProductImage(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	s, err := h.storeService.GetMyStore(ctx, userID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidProductID))
		return
	}

	var req dto.AddProductImageRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	image, err := h.productService.AddProductImage(ctx, productID.String(), s.ID, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Image added successfully", image)
}

// DeleteProductImage godoc
// @Summary      Delete Product Image
// @Description  Delete an image from a product
// @Tags         seller-products
// @Produce      json
// @Security     BearerAuth
// @Param 		 id   path 	string true "Product ID (UUID)" format(uuid)
// @Param 		 imageId path string true "Image ID (UUID)" format(uuid)
// @Success      200  {object}  dto.APIResponse
// @Failure      400  {object} 	dto.APIResponse
// @Router       /seller/products/{id}/images/{imageId} [delete]
func (h *ProductHandler) DeleteProductImage(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	s, err := h.storeService.GetMyStore(ctx, userID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidProductID))
		return
	}

	imageID, err := uuid.Parse(c.Param("imageId"))
	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, "invalid image id"))
		return
	}

	err = h.productService.DeleteProductImage(ctx, productID.String(), imageID.String(), s.ID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Image deleted successfully", nil)
}

// SetPrimaryImage godoc
// @Summary      Set Primary Product Image
// @Description  Set an image as the primary image for a product
// @Tags         seller-products
// @Produce      json
// @Security     BearerAuth
// @Param 		 id   path 	string true "Product ID (UUID)" format(uuid)
// @Param 		 imageId path string true "Image ID (UUID)" format(uuid)
// @Success      200  {object}  dto.APIResponse
// @Failure      400  {object} 	dto.APIResponse
// @Router       /seller/products/{id}/images/{imageId}/primary [put]
func (h *ProductHandler) SetPrimaryImage(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	s, err := h.storeService.GetMyStore(ctx, userID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidProductID))
		return
	}

	imageID, err := uuid.Parse(c.Param("imageId"))
	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, "invalid image id"))
		return
	}

	err = h.productService.SetPrimaryImage(ctx, productID.String(), imageID.String(), s.ID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Primary image updated successfully", nil)
}
