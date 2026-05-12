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

func (repo *userRepository) GetAll(ctx context.Context, filter *dto.GetUsersQuery) ([]entity.User, int, error) {
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

	argsData := append(args, filter.Limit, offset)

	rows, err := repo.db.QueryContext(ctx, dataQuery, argsData...)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err find user with query")
		return nil, 0, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Msg("failed to close rows")
		}
	}()
	users := []entity.User{}
	var total int

	for rows.Next() {

		var u entity.User
		var rowCount int

		err = rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Role,
			&u.IsActive,
			&u.CreatedAt,
			&u.UpdatedAt,
			&rowCount,
		)

		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Str("userID", u.ID).Msg("err mapping user row")
			return nil, 0, err
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
	argPos := 1

	// search
	if filter.Search != "" {
		query += fmt.Sprintf(" AND u.name ILIKE $%d", argPos)
		args = append(args, "%"+filter.Search+"%")
		argPos++
	}

	return query, args
}

func (repo *userRepository) Create(ctx context.Context, user *entity.User) error {
	_, err := repo.db.ExecContext(ctx, createUserQuery, user.ID, user.Name, user.Email, user.Password, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("err to create user")
		return err
	}
	return nil
}

func (repo *userRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	var u entity.User
	err := repo.db.QueryRowContext(ctx, getUserByIDQuery, id).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
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

func (repo *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var u entity.User
	err := repo.db.QueryRowContext(ctx, getUserByEmailQuery, email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("email", email).Msg("err get user by email")
		return nil, err
	}

	return &u, nil
}

func (repo *userRepository) Update(ctx context.Context, id string, user *entity.User) error {
	result, err := repo.db.ExecContext(ctx, updateUserQuery, user.Name, user.Email, user.Role, id)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err update user")
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err get rows affected")
		return err
	}
	if rows == 0 {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err not found update user")
		return dto.NewError(http.StatusNotFound, preference.ErrUserNotFound)
	}
	return nil
}

func (repo *userRepository) Delete(ctx context.Context, id string) error {
	result, err := repo.db.ExecContext(ctx, deleteUserQuery, id)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err delete user")
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err get rows affected")
		return err
	}

	if rows == 0 {
		zerolog.Ctx(ctx).Error().Err(err).Str("id", id).Msg("err not found delete user")
		return dto.NewError(http.StatusNotFound, preference.ErrUserNotFound)
	}

	return err

}
