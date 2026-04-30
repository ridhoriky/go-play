package services

import (
	"ne-project/src/internal/config/token"
	"ne-project/src/internal/repositories"
	"ne-project/src/internal/services/auth"
	"ne-project/src/internal/services/category"
	"ne-project/src/internal/services/product"
	"ne-project/src/internal/services/report"
	"ne-project/src/internal/services/transaction"
	"ne-project/src/internal/services/user"

	"github.com/jmoiron/sqlx"
)

type Services struct {
	Auth        auth.AuthServiceItf
	Category    category.CategoryServiceItf
	Product     product.ProductServiceItf
	Transaction transaction.TransactionServiceItf
	Report      report.ReportServiceItf
	User        user.UserServiceItf
}

func NewServices(repositories *repositories.Repositories, tokenSvc *token.Token, db *sqlx.DB) *Services {
	return &Services{
		Auth:        auth.NewAuthService(repositories.User, repositories.Auth, tokenSvc, db),
		Category:    category.NewCategoryService(repositories.Category),
		Product:     product.NewProductService(repositories.Product),
		Transaction: transaction.NewTransactionService(repositories.Transaction),
		Report:      report.NewReportService(repositories.Report),
		User:        user.NewUserService(repositories.User),
	}
}
