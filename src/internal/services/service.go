package services

import (
	"ne-project/src/internal/config/appconfig"
	"ne-project/src/internal/config/token"
	"ne-project/src/internal/repositories"
	"ne-project/src/internal/services/admin"
	"ne-project/src/internal/services/auth"
	"ne-project/src/internal/services/cart"
	"ne-project/src/internal/services/category"
	"ne-project/src/internal/services/order"
	"ne-project/src/internal/services/product"
	"ne-project/src/internal/services/report"
	"ne-project/src/internal/services/review"
	"ne-project/src/internal/services/store"
	"ne-project/src/internal/services/transaction"
	"ne-project/src/internal/services/user"
	"ne-project/src/internal/services/wishlist"

	"github.com/redis/go-redis/v9"
)

type Services struct {
	Auth         auth.AuthServiceItf
	Cart         cart.CartServiceItf
	Category     category.CategoryServiceItf
	Order        order.OrderServiceItf
	Product      product.ProductServiceItf
	Transaction  transaction.TransactionServiceItf
	Report       report.ReportServiceItf
	SellerReport report.SellerReportServiceItf
	Review       review.ReviewServiceItf
	User         user.UserServiceItf
	Store        store.StoreServiceItf
	Wishlist     wishlist.WishlistServiceItf
	Admin        admin.AdminServiceItf
	AdminReport  admin.AdminReportServiceItf
}

func NewServices(repositories *repositories.Repositories, tokenSvc *token.Token, rdb *redis.Client, cfg *appconfig.Config) *Services {
	return &Services{
		Auth:         auth.NewAuthService(repositories.User, repositories.Auth, tokenSvc),
		Cart:         cart.NewCartService(repositories.Cart, repositories.Product, repositories.Store, cfg),
		Category:     category.NewCategoryService(repositories.Category),
		Order:        order.NewOrderService(repositories.Order, repositories.Cart, repositories.Product, repositories.Store),
		Product:      product.NewProductService(repositories.Product, repositories.Wishlist, repositories.ProductImage, repositories.Order),
		Transaction:  transaction.NewTransactionService(repositories.Transaction),
		Report:       report.NewReportService(repositories.Report),
		SellerReport: report.NewSellerReportService(repositories.Report),
		Review:       review.NewReviewService(repositories.Review, repositories.Order, repositories.Product),
		User:         user.NewUserService(repositories.User, rdb),
		Store:        store.NewStoreService(repositories.Store, repositories.User),
		Wishlist:     wishlist.NewWishlistService(repositories.Wishlist, repositories.Product),
		Admin:        admin.NewAdminService(repositories.Admin),
		AdminReport:  admin.NewAdminReportService(repositories.Admin),
	}
}
