package auth

import (
	"context"
	"database/sql"
	"time"

	"ne-project/src/internal/models/entity"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

func (r *authRepository) SaveRefreshToken(ctx context.Context, tx *sqlx.Tx, tokenHash string, userID string, userAgent string, ipAddress string, expiresAt time.Time) error {
	query := createRefreshTokenQuery

	var exec sqlx.ExecerContext
	if tx != nil {
		exec = tx
	} else {
		exec = r.db
	}

	_, err := exec.ExecContext(ctx, query, tokenHash, userID, userAgent, ipAddress, expiresAt, time.Now())
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to save refresh token")
		return err
	}

	return nil
}

func (r *authRepository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	query := getRefreshTokenByHashQuery

	var token entity.RefreshToken
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(&token.ID, &token.TokenHash, &token.UserID, &token.UserAgent, &token.IPAddress, &token.ExpiresAt, &token.CreatedAt, &token.RevokedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to get refresh token")
		return nil, err
	}

	return &token, nil
}

func (r *authRepository) RevokeRefreshToken(ctx context.Context, tx *sqlx.Tx, tokenHash string) error {
	query := revokeRefreshTokenQuery

	var exec sqlx.ExecerContext
	if tx != nil {
		exec = tx
	} else {
		exec = r.db
	}

	result, err := exec.ExecContext(ctx, query, time.Now(), tokenHash)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to revoke refresh token")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		zerolog.Ctx(ctx).Warn().Str("token_hash", tokenHash).Msg("No refresh token found to revoke")
	}

	return nil
}
