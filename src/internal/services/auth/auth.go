package auth

// import (
// 	"context"
// 	"ne-project/src/internal/models/dto"
// 	"time"

// 	"github.com/rs/zerolog"
// )

// type AuthServiceItf interface {
// 	RequestChallenge(ctx context.Context, req *dto.RequestChallengeRequest) (*dto.ChallengeResponse, error)
// 	VerifyChallenge(ctx context.Context, req *dto.VerifyChallengeRequest) (*dto.ChallengeResponse, error)
// 	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
// 	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
// 	Logout(ctx context.Context, req *dto.LogoutRequest) (*dto.LogoutResponse, error)
// 	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
// }

// type auth struct {
// 	log            zerolog.Logger
// 	opt            MiddlewareOptions
// 	limit          int
// 	period         time.Duration
// 	userRepository UserRepositoryItf
// }

// func NewAuthService(log zerolog.Logger, opt MiddlewareOptions, limit int, period time.Duration, userRepository UserRepositoryItf) AuthServiceItf {
// 	return &auth{
// 		log:            log,
// 		opt:            opt,
// 		limit:          limit,
// 		period:         period,
// 		userRepository: userRepository,
// 	}
// }
