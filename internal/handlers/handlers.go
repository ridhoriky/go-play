package handlers

import (
	"ne-project/internal/handlers/category"
	"ne-project/internal/handlers/product"
	"ne-project/internal/services"

	"github.com/gin-gonic/gin"
)


type Handlers struct {
	Category category.CategoryHandlerItf
	Product  product.ProductHandlerItf
}

func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{
		Category: category.NewCategoryHandler(services.Category),
		Product:  product.NewProductHandler(services.Product),
	}
}

func (h *Handlers) RegisterRoutes(r *gin.Engine) {
	h.Category.RegisterRoutes(r)
	h.Product.RegisterRoutes(r)
}