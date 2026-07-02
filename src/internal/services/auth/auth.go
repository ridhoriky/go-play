package auth

import (
	"context"

	"ne-project/src/internal/config/token"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/repositories/auth"
	"ne-project/src/internal/repositories/user"
	"ne-project/src/internal/utils/mailer"
)

type AuthServiceItf interface {
	Login(ctx context.Context, email string, password string, userAgent string, ipAddr string) (*dto.LoginResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	RefreshToken(ctx context.Context, refreshToken string, userAgent string, ipAddr string) (*dto.AuthTokenResponse, error)
	Logout(ctx context.Context, userID string, refreshToken string) error
	VerifyEmail(ctx context.Context, email string, otpCode string) error
	ResendOTP(ctx context.Context, email string) error
	LoginGoogle(ctx context.Context, idToken string, userAgent string, ipAddr string) (*dto.LoginResponse, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token string, newPassword string) error
}

type authService struct {
	userRepository user.UserRepositoryItf
	authRepository auth.AuthRepositoryItf
	tokenService   *token.Token
	mailer         *mailer.Mailer
}

func NewAuthService(
	userRepository user.UserRepositoryItf,
	authRepository auth.AuthRepositoryItf,
	tokenService *token.Token,
	mailer *mailer.Mailer,
) AuthServiceItf {
	return &authService{
		userRepository: userRepository,
		authRepository: authRepository,
		tokenService:   tokenService,
		mailer:         mailer,
	}
}
