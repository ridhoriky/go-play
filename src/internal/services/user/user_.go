package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/utils/hash"
	"ne-project/src/internal/utils/validation"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

const (
	userCacheKey = "user:%s"
	userCacheTTL = 1 * time.Minute
)

func (s *userService) GetAllUsers(ctx context.Context, req *dto.GetUsersQuery) (*dto.UserListResponse, error) {
	req.Page = validation.ValidatePage(req.Page)
	req.Limit = validation.ValidatePageSize(req.Limit)

	users, total, err := s.userRepository.GetAll(ctx, req)
	if err != nil {
		return nil, err
	}

	res := make([]dto.UserResponse, 0, len(users))

	for i := range users {
		u := &users[i]
		resUser := dto.UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Role:      u.Role,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}

		res = append(res, resUser)
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	resp := &dto.UserListResponse{
		Data: res,
		Meta: dto.PaginationMeta{
			Total:      total,
			Page:       req.Page,
			Limit:      req.Limit,
			TotalPages: totalPages,
		},
	}

	return resp, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*entity.User, error) {
	key := fmt.Sprintf(userCacheKey, id)

	val, err := s.rdb.Get(ctx, key).Result()
	if err == nil {
		var u entity.User
		if err = json.Unmarshal([]byte(val), &u); err == nil {
			zerolog.Ctx(ctx).Debug().Str("id", id).Msg("user cache hit")
			return &u, nil
		}

	} else if !errors.Is(err, redis.Nil) {
		zerolog.Ctx(ctx).Warn().Err(err).Str("id", id).Msg("failed to get user from cache")
	}

	u, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(u)
	if err == nil {
		if err := s.rdb.Set(ctx, key, data, userCacheTTL).Err(); err != nil {
			zerolog.Ctx(ctx).Warn().Err(err).Str("id", id).Msg("failed to save user to cache")
		}
	}
	return u, nil
}

func (s *userService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*entity.User, error) {

	if req.Name == "" {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrUserNameRequired)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrUserNameRequired)
	}

	if req.Email == "" {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrUserEmailRequired)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrUserEmailRequired)
	}

	if req.Role == "" {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrUserRoleRequired)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrUserRoleRequired)
	}

	// Hash user password before persist
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to hash password")
		return nil, dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	u := &entity.User{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
		IsActive: true,
	}

	if err := s.userRepository.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userService) UpdateUser(ctx context.Context, id string, user *dto.UpdateUserRequest) (*entity.User, error) {
	existingUser, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.Name != nil && *user.Name == "" {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrUserNameRequired)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrUserNameRequired)
	}

	if user.Email != nil && *user.Email == "" {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrUserEmailRequired)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrUserEmailRequired)
	}

	if user.Role != nil && *user.Role == "" {
		zerolog.Ctx(ctx).Error().Msg(preference.ErrUserRoleRequired)
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrUserRoleRequired)
	}

	if user.Name != nil {
		existingUser.Name = *user.Name
	}
	if user.Email != nil {
		existingUser.Email = *user.Email
	}
	if user.Role != nil {
		existingUser.Role = *user.Role
	}
	if user.Phone != nil {
		existingUser.Phone = user.Phone
	}
	if user.AvatarURL != nil {
		existingUser.AvatarURL = user.AvatarURL
	}

	if err := s.userRepository.Update(ctx, id, existingUser); err != nil {
		return nil, err
	}
	return existingUser, nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepository.Delete(ctx, id)
}
