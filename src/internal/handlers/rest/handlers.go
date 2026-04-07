package rest

import (
	"ne-project/src/internal/handlers/rest/category"
	"ne-project/src/internal/handlers/rest/product"
	"ne-project/src/internal/handlers/rest/report"
	"ne-project/src/internal/handlers/rest/transaction"
	"ne-project/src/internal/services"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Category    category.CategoryHandlerItf
	Product     product.ProductHandlerItf
	Transaction transaction.TransactionHandlerItf
	Report      report.ReportHandlerItf
}

func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{
		Category:    category.NewCategoryHandler(services.Category),
		Product:     product.NewProductHandler(services.Product),
		Transaction: transaction.NewTransactionHandler(services.Transaction),
		Report:      report.NewReportHandler(services.Report),
	}
}

func (h *Handlers) RegisterRoutes(r *gin.RouterGroup) {
	h.Category.RegisterRoutes(r)
	h.Product.RegisterRoutes(r)
	h.Transaction.RegisterRoutes(r)
	h.Report.RegisterRoutes(r)
}
