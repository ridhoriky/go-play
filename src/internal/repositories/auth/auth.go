package auth

import (
	"context"
	"errors"
	"time"

	"ne-project/src/internal/models/entity"

	"github.com/redis/go-redis/v9"
)

type AuthRepositoryItf interface {
	SaveRefreshToken(ctx context.Context, tokenHash string, userID string, userAgent string, ipAddress string, expiresAt time.Time) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
}

type authRepository struct {
	rdb *redis.Client
}

type refreshTokenData struct {
	TokenHash string    `json:"token_hash"`
	UserID    string    `json:"user_id"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	refreshTokenPrefix = "refresh_token:"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

func NewAuthRepository(rdb *redis.Client) AuthRepositoryItf {
	return &authRepository{
		rdb: rdb,
	}
}
