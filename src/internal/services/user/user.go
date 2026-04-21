package user

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/repositories/user"
)

type UserServiceItf interface {
	GetAllUsers(ctx context.Context, query *dto.GetUsersQuery) (*dto.UserListResponse, error)
	GetUserByID(ctx context.Context, id string) (*entity.User, error)
	CreateUser(ctx context.Context, user *dto.CreateUserRequest) (*entity.User, error)
	UpdateUser(ctx context.Context, id string, user *dto.UpdateUserRequest) (*entity.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	userRepository user.UserRepositoryItf
}

func NewUserService(userRepository user.UserRepositoryItf) UserServiceItf {
	return &userService{
		userRepository: userRepository,
	}
}
