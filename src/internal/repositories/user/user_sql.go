package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"

	"github.com/rs/zerolog"
)

func (r *userRepository) GetAll(ctx context.Context, filter *dto.GetUsersQuery) ([]entity.User, int, error) {
	filterQuery, args := buildUserFilters(filter)
	dataQuery := getAllUsersQuery + filterQuery

	sortBy := "u.created_at"
	allowedSortColumns := map[string]string{
		"name":       "u.name",
		"created_at": "u.created_at",
	}

	if val, ok := allowedSortColumns[filter.SortBy]; ok {
		sortBy = val
	}

	sortDir := "ASC"
	if filter.SortDir == "DESC" {
		sortDir = "DESC"
	}

	dataQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortDir)

	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)

	offset := (filter.Page - 1) * filter.Limit

	args = append(args, filter.Limit, offset)

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err find user with query")
		return nil, 0, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			zerolog.Ctx(ctx).Error().Err(closeErr).Msg("failed to close rows")
		}
	}()
	users := []entity.User{}
	var total int

	for rows.Next() {

		var u entity.User
		var rowCount int

		if scanErr := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Role,
			&u.IsActive,
			&u.IsVerified,
			&u.CreatedAt,
			&u.UpdatedAt,
			&rowCount,
		); scanErr != nil {
			zerolog.Ctx(ctx).Error().Err(scanErr).Str("userID", u.ID).Msg("err mapping user row")
			return nil, 0, scanErr
		}

		users = append(users, u)

		// assign value total from firt row
		if total == 0 {
			total = rowCount
		}
	}

	if err = rows.Err(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err iterating user rows")
		return nil, 0, err
	}

	return users, total, nil
}

func buildUserFilters(filter *dto.GetUsersQuery) (string, []any) {

	query := ""
	args := []any{}

	if filter.Search != "" {
		query += " AND u.name ILIKE $1"
		args = append(args, "%"+filter.Search+"%")
	}

	return query, args
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	_, err := r.db.ExecContext(ctx, createUserQuery, user.ID, user.Name, user.Email, user.Password, user.Role, user.IsActive, user.IsVerified, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err to create user")
		return err
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	var u entity.User
	err := r.db.QueryRowContext(ctx, getUserByIDQuery, id).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.IsVerified, &u.Phone, &u.AvatarURL, &u.Address, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err user id not found")
		return nil, dto.NewError(http.StatusNotFound, preference.ErrUserNotFound)
	}

	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err get single user")
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var u entity.User
	err := r.db.QueryRowContext(ctx, getUserByEmailQuery, email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.IsVerified, &u.Phone, &u.AvatarURL, &u.Address, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sql.ErrNoRows
	}

	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("email", email).Msg("err get user by email")
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) Update(ctx context.Context, id string, user *entity.User) error {
	result, err := r.db.ExecContext(ctx, updateUserQuery, user.Name, user.Email, user.Password, user.Role, user.IsVerified, user.IsActive, user.Phone, user.AvatarURL, user.Address, id)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err update user")
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err get rows affected")
		return err
	}
	if rowsAffected == 0 {
		return dto.NewError(http.StatusNotFound, preference.ErrUserNotFound)
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, deleteUserQuery, id)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err delete user")
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err get rows affected")
		return err
	}

	if rowsAffected == 0 {
		return dto.NewError(http.StatusNotFound, preference.ErrUserNotFound)
	}

	return nil
}
