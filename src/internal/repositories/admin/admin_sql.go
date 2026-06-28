package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"ne-project/src/internal/models/dto"

	"github.com/jmoiron/sqlx"
)

type AdminRepository struct {
	db *sqlx.DB
}

func NewAdminRepository(db *sqlx.DB) AdminRepositoryItf {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) ListUsers(ctx context.Context, params *dto.AdminUserListParams) ([]dto.AdminUserResponse, int, error) {
	query := queryListUsersBase
	countQuery := queryCountUsersBase

	whereClauses := []string{}
	args := []any{}
	argId := 1

	if params.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(u.name ILIKE $%d OR u.email ILIKE $%d)", argId, argId))
		args = append(args, "%"+params.Search+"%")
		argId++
	}

	if params.Role != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("u.role = $%d", argId))
		args = append(args, params.Role)
		argId++
	}

	if params.IsActive != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.is_active = $%d", argId))
		args = append(args, *params.IsActive)
		argId++
	}

	if len(whereClauses) > 0 {
		query += " AND " + strings.Join(whereClauses, " AND ")
		countQuery += " AND " + strings.Join(whereClauses, " AND ")
	}

	// Default sort
	sortBy := "u.created_at"
	sortOrder := "DESC"

	if params.SortBy != "" {
		switch params.SortBy {
		case "name":
			sortBy = "u.name"
		case "email":
			sortBy = "u.email"
		case "created_at":
			sortBy = "u.created_at"
		}
	}
	if strings.ToUpper(params.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", sortBy, sortOrder, argId, argId+1)

	limit := params.Limit
	if limit == 0 {
		limit = 10
	}
	offset := max(0, (params.Page-1)*limit)

	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)

	var users []dto.AdminUserResponse
	err = r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *AdminRepository) GetUserByID(ctx context.Context, id string) (*dto.AdminUserResponse, error) {
	var user dto.AdminUserResponse
	err := r.db.GetContext(ctx, &user, queryGetUserByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err

	}
	return &user, nil
}

func (r *AdminRepository) UpdateUser(ctx context.Context, id string, req dto.AdminUpdateUserRequest) error {
	result, err := r.db.ExecContext(ctx, queryUpdateUser, id, req.Role, req.IsActive)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *AdminRepository) ListSellers(ctx context.Context, params *dto.AdminSellerListParams) ([]dto.AdminSellerResponse, int, error) {
	query := queryListSellersBase
	countQuery := queryCountSellersBase

	whereClauses := []string{}
	args := []any{}
	argId := 1

	if params.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(s.store_name ILIKE $%d OR u.name ILIKE $%d)", argId, argId))
		args = append(args, "%"+params.Search+"%")
		argId++
	}

	if params.IsVerified != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("s.is_verified = $%d", argId))
		args = append(args, *params.IsVerified)
		argId++
	}

	if len(whereClauses) > 0 {
		query += " AND " + strings.Join(whereClauses, " AND ")
		countQuery += " AND " + strings.Join(whereClauses, " AND ")
	}

	sortBy := "s.created_at"
	sortOrder := "DESC"

	if params.SortBy != "" {
		switch params.SortBy {
		case "store_name":
			sortBy = "s.store_name"
		case "rating_avg":
			sortBy = "s.rating_avg"
		case "total_sales":
			sortBy = "s.total_sales"
		case "created_at":
			sortBy = "s.created_at"
		}
	}
	if strings.ToUpper(params.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", sortBy, sortOrder, argId, argId+1)

	limit := params.Limit
	if limit == 0 {
		limit = 10
	}
	offset := max(0, (params.Page-1)*limit)

	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)

	var sellers []dto.AdminSellerResponse
	err = r.db.SelectContext(ctx, &sellers, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return sellers, total, nil
}

func (r *AdminRepository) VerifySeller(ctx context.Context, storeID string) error {
	result, err := r.db.ExecContext(ctx, queryUpdateStoreVerification, storeID, true)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("store not found")
	}
	return nil
}

func (r *AdminRepository) UnverifySeller(ctx context.Context, storeID string) error {
	result, err := r.db.ExecContext(ctx, queryUpdateStoreVerification, storeID, false)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("store not found")
	}
	return nil
}

func (r *AdminRepository) GetPlatformSummary(ctx context.Context) (*dto.AdminPlatformSummary, error) {
	var summary dto.AdminPlatformSummary
	err := r.db.GetContext(ctx, &summary, queryGetPlatformSummary)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

func (r *AdminRepository) GetTopStores(ctx context.Context, sortBy string, limit int) ([]dto.AdminSellerResponse, error) {
	query := queryGetTopStoresBase

	orderBy := "revenue DESC"
	switch sortBy {
	case "revenue":
		orderBy = "revenue DESC"
	case "total_sales":
		orderBy = "s.total_sales DESC"
	case "rating":
		orderBy = "s.rating_avg DESC"
	case "products":
		orderBy = "total_products DESC"
	}

	query += fmt.Sprintf(" ORDER BY %s LIMIT $1", orderBy)

	var stores []dto.AdminSellerResponse
	err := r.db.SelectContext(ctx, &stores, query, limit)
	if err != nil {
		return nil, err
	}

	return stores, nil
}

func (r *AdminRepository) GetTopProducts(ctx context.Context, sortBy string, limit int) ([]dto.AdminTopProductResponse, error) {
	query := queryGetTopProductsBase

	orderBy := "revenue DESC"
	switch sortBy {
	case "revenue":
		orderBy = "revenue DESC"
	case "quantity":
		orderBy = "quantity_sold DESC"
	case "rating":
		orderBy = "p.rating_avg DESC"
	}

	query += fmt.Sprintf(" ORDER BY %s LIMIT $1", orderBy)

	var products []dto.AdminTopProductResponse
	err := r.db.SelectContext(ctx, &products, query, limit)
	if err != nil {
		return nil, err
	}

	return products, nil
}
