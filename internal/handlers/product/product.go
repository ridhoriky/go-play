package product

import (
	"ne-project/internal/services/product"

	"github.com/gin-gonic/gin"
)


type ProductHandlerItf interface {
	RegisterRoutes(r *gin.Engine)
}

type productHandler struct {
	productService product.ProductServiceItf
}

func NewProductHandler(productService product.ProductServiceItf) ProductHandlerItf {
	return &productHandler{
		productService: productService,
	}
}

func (h *productHandler) RegisterRoutes(r *gin.Engine) {
	productRoutes := r.Group("/products")
	{
		productRoutes.GET("", h.GetAll)
		productRoutes.POST("", h.Create)
		productRoutes.GET("/:id", h.GetByID)
		productRoutes.PUT("/:id", h.Update)
		productRoutes.DELETE("/:id", h.Delete)
		productRoutes.POST("/bulk", h.CreateMultiple)
	}
}