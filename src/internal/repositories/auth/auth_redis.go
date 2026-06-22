package auth

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"ne-project/src/internal/models/entity"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func (r *authRepository) key(tokenHash string) string {
	return refreshTokenPrefix + tokenHash
}

func (r *authRepository) SaveRefreshToken(ctx context.Context, tokenHash string, userID string, userAgent string, ipAddress string, expiresAt time.Time) error {
	data := refreshTokenData{
		TokenHash: tokenHash,
		UserID:    userID,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to marshal refresh token data")
		return err
	}

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return errors.New("token already expired")
	}

	err = r.rdb.Set(ctx, r.key(tokenHash), jsonData, ttl).Err()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to save refresh token to Redis")
		return err
	}

	return nil
}

func (r *authRepository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	jsonData, err := r.rdb.Get(ctx, r.key(tokenHash)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrRefreshTokenNotFound
		}
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to get refresh token from Redis")
		return nil, err
	}

	var data refreshTokenData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to unmarshal refresh token data")
		return nil, err
	}

	return &entity.RefreshToken{
		ID:        tokenHash,
		UserID:    data.UserID,
		TokenHash: data.TokenHash,
		UserAgent: data.UserAgent,
		IPAddress: data.IPAddress,
		ExpiresAt: data.ExpiresAt,
		CreatedAt: data.CreatedAt,
		RevokedAt: nil,
	}, nil
}

func (r *authRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	result, err := r.rdb.Del(ctx, r.key(tokenHash)).Result()
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to revoke refresh token from Redis")
		return err
	}

	if result == 0 {
		zerolog.Ctx(ctx).Warn().Str("token_hash", tokenHash).Msg("No refresh token found to revoke")
	}

	return nil
}
