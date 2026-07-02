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
	Auth          rest.AuthHandler
	Category      rest.CategoryHandler
	Product       rest.ProductHandler
	Transaction   rest.TransactionHandler
	Report        rest.ReportHandler
	SellerReport  rest.SellerReportHandler
	User          rest.UserHandler
	System        rest.SystemHandler
	Store         rest.StoreHandler
	Cart          rest.CartHandler
	Order         rest.OrderHandler
	Review        rest.ReviewHandler
	Wishlist      rest.WishlistHandler
	SellerProfile rest.SellerProfileHandler
	Admin         rest.AdminHandler
	AdminReport   rest.AdminReportHandler
	Payment       rest.PaymentHandler
}

func NewHandlers(db *sqlx.DB, services *services.Services) *Handlers {
	return &Handlers{
		Auth:          *rest.NewAuthHandler(services.Auth),
		Category:      *rest.NewCategoryHandler(services.Category),
		Product:       *rest.NewProductHandler(services.Product, services.Store),
		Transaction:   *rest.NewTransactionHandler(services.Transaction),
		Report:        *rest.NewReportHandler(services.Report),
		SellerReport:  *rest.NewSellerReportHandler(services.SellerReport),
		User:          *rest.NewUserHandler(services.User),
		System:        *rest.NewSystemHandler(db),
		Store:         *rest.NewStoreHandler(services.Store),
		Cart:          *rest.NewCartHandler(services.Cart),
		Order:         *rest.NewOrderHandler(services.Order),
		Review:        *rest.NewReviewHandler(services.Review),
		Wishlist:      *rest.NewWishlistHandler(services.Wishlist),
		SellerProfile: *rest.NewSellerProfileHandler(services.User, services.Store),
		Admin:         *rest.NewAdminHandler(services.Admin),
		AdminReport:   *rest.NewAdminReportHandler(services.AdminReport),
		Payment:       *rest.NewPaymentHandler(services.Payment),
	}
}

func (h *Handlers) RegisterRoutes(r *gin.RouterGroup, tokenSvc *token.Token, mw *middleware.Middleware) {
	h.registerPublicRoutes(r, tokenSvc, mw)
	h.registerProtectedRoutes(r, tokenSvc, mw)
}

func (h *Handlers) registerPublicRoutes(r *gin.RouterGroup, tokenSvc *token.Token, mw *middleware.Middleware) {
	// SYSTEM ROUTES
	systemGroup := r.Group("/system")
	systemGroup.Use(mw.SetComponent("system"))
	{
		systemGroup.GET("/health", h.System.HealthCheck)
		systemGroup.GET("/ready", h.System.ReadyCheck)
	}

	// PUBLIC AUTH ROUTES
	authGroup := r.Group("/auth")
	authGroup.Use(mw.RateLimiter(), mw.SetComponent("auth"))
	{
		authGroup.POST("/login", h.Auth.Login)
		authGroup.POST("/register", h.Auth.Register)
		authGroup.POST("/refresh", h.Auth.RefreshToken)
		authGroup.POST("/verify-email", h.Auth.VerifyEmail)
		authGroup.POST("/resend-otp", h.Auth.ResendOTP)
		authGroup.POST("/google", h.Auth.GoogleLogin)
		authGroup.POST("/forgot-password", h.Auth.ForgotPassword)
		authGroup.POST("/reset-password", h.Auth.ResetPassword)
	}

	// PUBLIC STORE ROUTES
	publicStoreGroup := r.Group("/stores")
	publicStoreGroup.Use(mw.RateLimiter(), mw.SetComponent("store"), mw.OptionalJWTAuth(tokenSvc))
	{
		publicStoreGroup.GET("", h.Store.ListStores)
		publicStoreGroup.GET("/:slug", h.Store.GetStoreBySlug)
		publicStoreGroup.GET("/:slug/products", h.Product.GetStoreProducts)
	}

	// PUBLIC PRODUCT ROUTES
	publicProductGroup := r.Group("/products")
	publicProductGroup.Use(mw.RateLimiter(), mw.SetComponent("product"), mw.OptionalJWTAuth(tokenSvc))
	{
		publicProductGroup.GET("", h.Product.GetAll)
		publicProductGroup.GET("/:id", h.Product.GetByID)
		publicProductGroup.GET("/:id/reviews", h.Review.GetProductReviews)
	}

	// PUBLIC CATEGORY ROUTES
	publicCategoryGroup := r.Group("/categories")
	publicCategoryGroup.Use(mw.RateLimiter(), mw.SetComponent("category"))
	{
		publicCategoryGroup.GET("", h.Category.GetAll)
		publicCategoryGroup.GET("/:id", h.Category.GetByID)
	}

	// PUBLIC PAYMENT ROUTES
	publicPaymentGroup := r.Group("/payments")
	publicPaymentGroup.Use(mw.RateLimiter(), mw.SetComponent("payment"))
	{
		publicPaymentGroup.POST("/callback", h.Payment.PaymentCallback)
		publicPaymentGroup.GET("/callback", h.Payment.PaymentCallback) // Supported for simulated redirects
	}
}

func (h *Handlers) registerProtectedRoutes(r *gin.RouterGroup, tokenSvc *token.Token, mw *middleware.Middleware) {
	protected := r.Group("")
	protected.Use(mw.JWTAuth(tokenSvc), mw.RateLimiter())

	// Logout endpoint (protected)
	protected.POST("/auth/logout", h.Auth.Logout)

	h.registerCommonProtectedRoutes(protected, mw)
	h.registerSellerRoutes(protected, mw)
	h.registerBuyerRoutes(protected, mw)
	h.registerAdminRoutes(protected, mw)
}

func (h *Handlers) registerCommonProtectedRoutes(protected *gin.RouterGroup, mw *middleware.Middleware) {
	// products
	routerProducts := protected.Group("/products")
	routerProducts.Use(mw.SetComponent("product"))
	{
		routerProducts.POST("", mw.RequireRole("seller", "admin"), h.Product.Create)
		routerProducts.PUT("/:id", mw.RequireRole("seller", "admin"), h.Product.Update)
		routerProducts.DELETE("/:id", mw.RequireRole("seller", "admin"), h.Product.Delete)
		routerProducts.POST("/bulk", mw.RequireRole("seller", "admin"), h.Product.CreateMultiple)
	}

	// transactions
	routerTransactions := protected.Group("/transactions")
	routerTransactions.Use(mw.SetComponent("transaction"))
	{
		routerTransactions.POST("", mw.RequireRole("buyer", "admin"), h.Transaction.Checkout)
		routerTransactions.GET("/:id", h.Transaction.GetByID)
		routerTransactions.PATCH("/:id", mw.RequireRole("seller", "admin"), h.Transaction.UpdateTransactionStatus)
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
}

func (h *Handlers) registerSellerRoutes(protected *gin.RouterGroup, mw *middleware.Middleware) {
	// stores (protected seller routes)
	routerSellerStores := protected.Group("/seller/stores")
	routerSellerStores.Use(mw.SetComponent("store"))
	{
		routerSellerStores.POST("", mw.RequireRole("buyer", "seller", "admin"), h.Store.CreateStore)
		routerSellerStores.GET("/me", mw.RequireRole("seller", "admin"), h.Store.GetMyStore)
		routerSellerStores.PUT("/me", mw.RequireRole("seller", "admin"), h.Store.UpdateStore)
	}

	// seller products
	routerSellerProducts := protected.Group("/seller/products")
	routerSellerProducts.Use(mw.SetComponent("product"), mw.RequireRole("seller", "admin"))
	{
		routerSellerProducts.GET("", h.Product.GetMyProducts)
		routerSellerProducts.GET("/:id", h.Product.GetSellerProductDetail)
		routerSellerProducts.PUT("/:id", h.Product.UpdateSellerProduct)
		routerSellerProducts.DELETE("/:id", h.Product.DeleteSellerProduct)
		routerSellerProducts.POST("/:id/images", h.Product.AddProductImage)
		routerSellerProducts.DELETE("/:id/images/:imageId", h.Product.DeleteProductImage)
		routerSellerProducts.PUT("/:id/images/:imageId/primary", h.Product.SetPrimaryImage)
	}

	sellerGroup := protected.Group("/seller")
	sellerGroup.Use(mw.RequireRole("seller"))
	{
		sellerGroup.GET("/profile", h.SellerProfile.GetSellerProfile)
		sellerGroup.PUT("/profile", h.SellerProfile.UpdateSellerProfile)
		sellerGroup.GET("/store/stats", h.SellerProfile.GetStoreStats)

		sellerGroup.GET("/orders", h.Order.GetSellerOrders)
		sellerGroup.GET("/orders/:id", h.Order.GetSellerOrderDetail)
		sellerGroup.PATCH("/orders/:id/status", h.Order.UpdateSellerOrderStatus)

		sellerGroup.GET("/reports/summary", h.SellerReport.GetSalesSummary)
		sellerGroup.GET("/reports/top-products", h.SellerReport.GetTopProducts)
		sellerGroup.GET("/reports/dashboard", h.SellerReport.GetDashboard)

		sellerGroup.PUT("/reviews/:id/reply", h.Review.ReplyToReview)
	}
}

func (h *Handlers) registerBuyerRoutes(protected *gin.RouterGroup, mw *middleware.Middleware) {
	// Buyer group
	buyerGroup := protected.Group("")
	buyerGroup.Use(mw.RequireRole("buyer", "seller"))
	{
		// Cart Routes
		buyerGroup.POST("/cart", h.Cart.AddToCart)
		buyerGroup.GET("/cart", h.Cart.GetCart)
		buyerGroup.PUT("/cart/:id", h.Cart.UpdateCartItem)
		buyerGroup.DELETE("/cart/:id", h.Cart.RemoveFromCart)
		buyerGroup.DELETE("/cart", h.Cart.ClearCart)

		buyerGroup.POST("/orders", h.Order.Checkout)
		buyerGroup.GET("/orders", h.Order.GetBuyerOrders)
		buyerGroup.GET("/orders/:id", h.Order.GetBuyerOrderDetail)
		buyerGroup.PATCH("/orders/:id/cancel", h.Order.CancelOrder)
		buyerGroup.PATCH("/orders/:id/confirm", h.Order.ConfirmReceived)
		buyerGroup.POST("/orders/:id/pay", h.Payment.CreatePayment)
		buyerGroup.GET("/orders/:id/payment-status", h.Payment.GetPaymentStatus)

		// Review Routes
		buyerGroup.POST("/reviews", h.Review.CreateReview)

		// Wishlist Routes
		buyerGroup.GET("/wishlist", h.Wishlist.GetWishlist)
		buyerGroup.POST("/wishlist", h.Wishlist.AddToWishlist)
		buyerGroup.POST("/wishlist/toggle", h.Wishlist.ToggleWishlist)
		buyerGroup.DELETE("/wishlist/:productId", h.Wishlist.RemoveFromWishlist)
	}
}

func (h *Handlers) registerAdminRoutes(protected *gin.RouterGroup, mw *middleware.Middleware) {
	// Admin group (placeholder for future routes)
	adminGroup := protected.Group("/admin")
	adminGroup.Use(mw.RequireRole("admin"))
	{
		// Users
		adminGroup.GET("/users", h.Admin.ListUsers)
		adminGroup.GET("/users/:id", h.Admin.GetUserDetail)
		adminGroup.PUT("/users/:id", h.Admin.UpdateUser)

		// Sellers
		adminGroup.GET("/sellers", h.Admin.ListSellers)
		adminGroup.PATCH("/sellers/:id/verify", h.Admin.VerifySeller)
		adminGroup.PATCH("/sellers/:id/unverify", h.Admin.UnverifySeller)

		// Categories
		adminGroup.POST("/categories", h.Category.Create)
		adminGroup.PUT("/categories/:id", h.Category.Update)
		adminGroup.DELETE("/categories/:id", h.Category.Delete)
		adminGroup.GET("/categories/tree", h.Category.GetCategoryTree)

		// Reports
		adminGroup.GET("/reports/summary", h.AdminReport.GetPlatformSummary)
		adminGroup.GET("/reports/top-stores", h.AdminReport.GetTopStores)
		adminGroup.GET("/reports/top-products", h.AdminReport.GetTopProducts)
	}
}
