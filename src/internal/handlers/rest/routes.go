package rest

import (
	"time"

	"ne-project/src/internal/config/middleware"
	"ne-project/src/internal/config/token"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) RegisterRoutes(r *gin.RouterGroup, tokenSvc token.Token, mw middleware.Middleware) {

	// ─── PUBLIC AUTH ROUTES  ────────────────────
	authGroup := r.Group("/auth")
	authGroup.Use(mw.RateLimiter(5, time.Minute))
	{
		authGroup.POST("/login", h.Auth.Login)
		authGroup.POST("/register", h.Auth.Register)
		authGroup.POST("/refresh", h.Auth.RefreshToken)
	}

	// ─── PROTECTED ROUTES  ────────────────────
	protected := r.Group("")
	protected.Use(mw.JWTAuth(tokenSvc))
	{
		// category
		routerCategories := protected.Group("/categories")
		{
			routerCategories.GET("", h.Category.GetAll)
			routerCategories.POST("", h.Category.Create)
			routerCategories.GET("/:id", h.Category.GetByID)
			routerCategories.PUT("/:id", h.Category.Update)
			routerCategories.DELETE("/:id", h.Category.Delete)
		}

		// products
		routerProducts := protected.Group("/products")
		{
			routerProducts.GET("", h.Product.GetAll)
			routerProducts.POST("", h.Product.Create)
			routerProducts.GET("/:id", h.Product.GetByID)
			routerProducts.PUT("/:id", h.Product.Update)
			routerProducts.DELETE("/:id", h.Product.Delete)
			routerProducts.POST("/bulk", h.Product.CreateMultiple)
		}

		// transactions
		routerTransactions := protected.Group("/transactions")
		{
			routerTransactions.POST("", h.Transaction.Checkout)
			routerTransactions.GET("/:id", h.Transaction.GetByID)
			routerTransactions.PATCH("/:id", h.Transaction.UpdateTransactionStatus)
		}

		// reports
		routerReports := protected.Group("/reports")
		{
			routerReports.GET("/summary", h.Report.GetReports)
			routerReports.GET("/top-products", h.Report.GetTopProducts)
		}

		// users
		routerUsers := protected.Group("/users")
		{
			routerUsers.GET("", h.User.GetAllUser)
			routerUsers.POST("", h.User.CreateUser)
			routerUsers.GET("/:id", h.User.GetUserByID)
			routerUsers.PUT("/:id", h.User.UpdateUser)
			routerUsers.DELETE("/:id", h.User.DeleteUser)
		}

		// Logout endpoint (protected)
		protected.POST("/auth/logout", h.Auth.Logout)
	}
}
