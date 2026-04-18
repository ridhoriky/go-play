package rest

import "github.com/gin-gonic/gin"

func (h *Handlers) RegisterRoutes(r *gin.RouterGroup) {

	//category
	routerCategories := r.Group("/categories")
	routerCategories.GET("", h.Category.GetAll)
	routerCategories.POST("", h.Category.Create)
	routerCategories.GET("/:id", h.Category.GetByID)
	routerCategories.PUT("/:id", h.Category.Update)
	routerCategories.DELETE("/:id", h.Category.Delete)

	//products
	routerProducts := r.Group("/products")
	routerProducts.GET("", h.Product.GetAll)
	routerProducts.POST("", h.Product.Create)
	routerProducts.GET("/:id", h.Product.GetByID)
	routerProducts.PUT("/:id", h.Product.Update)
	routerProducts.DELETE("/:id", h.Product.Delete)
	routerProducts.POST("/bulk", h.Product.CreateMultiple)

	//transactions
	routerTransactions := r.Group("/transactions")
	routerTransactions.POST("", h.Transaction.Checkout)
	routerTransactions.GET("/:id", h.Transaction.GetByID)
	routerTransactions.PATCH("/:id", h.Transaction.UpdateTransactionStatus)

	//reports
	routerReports := r.Group("/reports")
	routerReports.GET("/summary", h.Report.GetReports)
	routerReports.GET("/top-products", h.Report.GetTopProducts)

}
