package repositories

import (
	"ne-project/src/internal/repositories/admin"
	"ne-project/src/internal/repositories/auth"
	"ne-project/src/internal/repositories/cart"
	"ne-project/src/internal/repositories/category"
	"ne-project/src/internal/repositories/order"
	"ne-project/src/internal/repositories/product"
	"ne-project/src/internal/repositories/product_image"
	"ne-project/src/internal/repositories/report"
	"ne-project/src/internal/repositories/review"
	"ne-project/src/internal/repositories/store"
	"ne-project/src/internal/repositories/transaction"
	"ne-project/src/internal/repositories/user"
	"ne-project/src/internal/repositories/wishlist"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	Auth         auth.AuthRepositoryItf
	Cart         cart.CartRepositoryItf
	Category     category.CategoryRepositoryItf
	Order        order.OrderRepositoryItf
	Product      product.ProductRepositoryItf
	Transaction  transaction.TransactionRepositoryItf
	Report       report.ReportRepositoryItf
	Review       review.ReviewRepositoryItf
	User         user.UserRepositoryItf
	Store        store.StoreRepositoryItf
	Wishlist     wishlist.WishlistRepositoryItf
	ProductImage product_image.ProductImageRepositoryItf
	Admin        admin.AdminRepositoryItf
}

func NewRepository(db *sqlx.DB, rdb *redis.Client) *Repositories {
	return &Repositories{
		Auth:         auth.NewAuthRepository(rdb),
		Cart:         cart.NewCartRepository(db),
		Category:     category.NewCategoryRepository(db),
		Order:        order.NewOrderRepository(db),
		Product:      product.NewProductRepository(db),
		Transaction:  transaction.NewTransactionRepository(db),
		Report:       report.NewReportRepository(db),
		Review:       review.NewReviewRepository(db),
		User:         user.NewUserRepository(db),
		Store:        store.NewStoreRepository(db),
		Wishlist:     wishlist.NewWishlistRepository(db),
		ProductImage: product_image.NewProductImageRepository(db),
		Admin:        admin.NewAdminRepository(db),
	}
}
