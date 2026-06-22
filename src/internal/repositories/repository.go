package repositories

import (
	"ne-project/src/internal/repositories/auth"
	"ne-project/src/internal/repositories/category"
	"ne-project/src/internal/repositories/product"
	"ne-project/src/internal/repositories/report"
	"ne-project/src/internal/repositories/transaction"
	"ne-project/src/internal/repositories/user"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	Auth        auth.AuthRepositoryItf
	Category    category.CategoryRepositoryItf
	Product     product.ProductRepositoryItf
	Transaction transaction.TransactionRepositoryItf
	Report      report.ReportRepositoryItf
	User        user.UserRepositoryItf
}

func NewRepository(db *sqlx.DB, rdb *redis.Client) *Repositories {
	return &Repositories{
		Auth:        auth.NewAuthRepository(rdb),
		Category:    category.NewCategoryRepository(db),
		Product:     product.NewProductRepository(db),
		Transaction: transaction.NewTransactionRepository(db),
		Report:      report.NewReportRepository(db),
		User:        user.NewUserRepository(db),
	}
}
