package user

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"

	"github.com/jmoiron/sqlx"
)

type UserRepositoryItf interface {
	GetAll(ctx context.Context, query *dto.GetUsersQuery) ([]entity.User, int, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, id string, user *entity.User) error
	Delete(ctx context.Context, id string) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepositoryItf {
	return &userRepository{
		db: db,
	}
}
