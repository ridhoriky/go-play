package routes

import (
	"ne-project/src/internal/config/token"
	"ne-project/src/internal/handlers/rest"
	"ne-project/src/internal/handlers/rest/middleware"

	"ne-project/src/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Handlers struct {
	Auth        rest.AuthHandler
	Category    rest.CategoryHandler
	Product     rest.ProductHandler
	Transaction rest.TransactionHandler
	Report      rest.ReportHandler
	User        rest.UserHandler
	System      rest.SystemHandler
}

func NewHandlers(db *sqlx.DB, services *services.Services) *Handlers {
	return &Handlers{
		Auth:        *rest.NewAuthHandler(services.Auth),
		Category:    *rest.NewCategoryHandler(services.Category),
		Product:     *rest.NewProductHandler(services.Product),
		Transaction: *rest.NewTransactionHandler(services.Transaction),
		Report:      *rest.NewReportHandler(services.Report),
		User:        *rest.NewUserHandler(services.User),
		System:      *rest.NewSystemHandler(db),
	}
}

func (h *Handlers) RegisterRoutes(r *gin.RouterGroup, tokenSvc *token.Token, mw middleware.Middleware) {

	// ─── SYSTEM ROUTES ──────────────────────────
	systemGroup := r.Group("/system")
	systemGroup.Use(mw.SetComponent("system"))
	{
		systemGroup.GET("/health", h.System.HealthCheck)
		systemGroup.GET("/ready", h.System.ReadyCheck)
	}

	// ─── PUBLIC AUTH ROUTES  ────────────────────
	authGroup := r.Group("/auth")
	authGroup.Use(mw.RateLimiter(), mw.SetComponent("auth"))
	{
		authGroup.POST("/login", h.Auth.Login)
		authGroup.POST("/register", h.Auth.Register)
		authGroup.POST("/refresh", h.Auth.RefreshToken)
	}

	// ─── PROTECTED ROUTES  ────────────────────
	protected := r.Group("")
	protected.Use(mw.JWTAuth(tokenSvc), mw.RateLimiter())
	{
		// category
		routerCategories := protected.Group("/categories")
		routerCategories.Use(mw.SetComponent("category"))
		{
			routerCategories.GET("", h.Category.GetAll)
			routerCategories.POST("", h.Category.Create)
			routerCategories.GET("/:id", h.Category.GetByID)
			routerCategories.PUT("/:id", h.Category.Update)
			routerCategories.DELETE("/:id", h.Category.Delete)
		}

		// products
		routerProducts := protected.Group("/products")
		routerProducts.Use(mw.SetComponent("product"))
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
		routerTransactions.Use(mw.SetComponent("transaction"))
		{
			routerTransactions.POST("", h.Transaction.Checkout)
			routerTransactions.GET("/:id", h.Transaction.GetByID)
			routerTransactions.PATCH("/:id", h.Transaction.UpdateTransactionStatus)
		}

		// reports
		routerReports := protected.Group("/reports")
		routerReports.Use(mw.SetComponent("report"))
		{
			routerReports.GET("/summary", h.Report.GetReports)
			routerReports.GET("/top-products", h.Report.GetTopProducts)
		}

		// users
		routerUsers := protected.Group("/users")
		routerUsers.Use(mw.SetComponent("user"))
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
