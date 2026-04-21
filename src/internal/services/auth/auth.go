package auth

import (
	"context"

	"ne-project/src/internal/config/token"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/repositories/auth"
	"ne-project/src/internal/repositories/user"

	"github.com/jmoiron/sqlx"
)

type AuthServiceItf interface {
	Login(ctx context.Context, email string, password string, userAgent string, ipAddr string) (*dto.LoginResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	RefreshToken(ctx context.Context, refreshToken string, userAgent string, ipAddr string) (*dto.AuthTokenResponse, error)
	Logout(ctx context.Context, userID string, refreshToken string) error
}

type authService struct {
	userRepository user.UserRepositoryItf
	authRepository auth.AuthRepositoryItf
	tokenService   token.Token
	db             *sqlx.DB
}

func NewAuthService(userRepository user.UserRepositoryItf, authRepository auth.AuthRepositoryItf, tokenService token.Token, db *sqlx.DB) AuthServiceItf {
	return &authService{
		userRepository: userRepository,
		authRepository: authRepository,
		tokenService:   tokenService,
		db:             db,
	}
}
