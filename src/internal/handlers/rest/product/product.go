package product

import (
	"ne-project/src/internal/services/product"

	"github.com/gin-gonic/gin"
)

type ProductHandlerItf interface {
	GetAll(c *gin.Context)
	Create(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	CreateMultiple(c *gin.Context)
}

type productHandler struct {
	productService product.ProductServiceItf
}

func NewProductHandler(productService product.ProductServiceItf) ProductHandlerItf {
	return &productHandler{
		productService: productService,
	}
}
