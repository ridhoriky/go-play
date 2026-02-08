package handlers

import (
	"ne-project/internal/handlers/category"
	"ne-project/internal/handlers/product"
	"ne-project/internal/handlers/transaction"
	"ne-project/internal/services"

	"github.com/gin-gonic/gin"
)


type Handlers struct {
	Category category.CategoryHandlerItf
	Product  product.ProductHandlerItf
	Transaction  transaction.TransactionHandlerItf
}

func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{
		Category: category.NewCategoryHandler(services.Category),
		Product:  product.NewProductHandler(services.Product),
		Transaction : transaction.NewTransactionHandler(services.Transaction),
	}
}

func (h *Handlers) RegisterRoutes(r *gin.Engine) {
	h.Category.RegisterRoutes(r)
	h.Product.RegisterRoutes(r)
	h.Transaction.RegisterRoutes(r)
}