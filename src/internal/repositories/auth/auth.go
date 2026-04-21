package auth

import (
	"context"
	"time"

	"ne-project/src/internal/models/entity"

	"github.com/jmoiron/sqlx"
)

type AuthRepositoryItf interface {
	SaveRefreshToken(ctx context.Context, tx *sqlx.Tx, tokenHash string, userID string, userAgent string, ipAddress string, expiresAt time.Time) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tx *sqlx.Tx, tokenHash string) error
}

type authRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) AuthRepositoryItf {
	return &authRepository{
		db: db,
	}
}
