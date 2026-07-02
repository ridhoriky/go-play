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

type otpData struct {
	Code     string `json:"code"`
	Attempts int    `json:"attempts"`
}

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

	// Track token hash for user session invalidation
	userTokensKey := "user_tokens:" + userID
	if err = r.rdb.SAdd(ctx, userTokensKey, tokenHash).Err(); err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Str("user_id", userID).Msg("Failed to add refresh token to user set")
	} else {
		_ = r.rdb.Expire(ctx, userTokensKey, 30*24*time.Hour)
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
	storedToken, err := r.GetRefreshTokenByHash(ctx, tokenHash)
	if err == nil && storedToken != nil {
		userTokensKey := "user_tokens:" + storedToken.UserID
		_ = r.rdb.SRem(ctx, userTokensKey, tokenHash).Err()
	}

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

func (r *authRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	userTokensKey := "user_tokens:" + userID
	hashes, err := r.rdb.SMembers(ctx, userTokensKey).Result()
	if err != nil {
		return err
	}

	for _, tokenHash := range hashes {
		_ = r.rdb.Del(ctx, r.key(tokenHash)).Err()
	}

	return r.rdb.Del(ctx, userTokensKey).Err()
}

func (r *authRepository) otpKey(email string) string {
	return "otp:code:" + email
}

func (r *authRepository) cooldownKey(email string) string {
	return "otp:cooldown:" + email
}

func (r *authRepository) resetTokenKey(token string) string {
	return "password_reset:token:" + token
}

func (r *authRepository) SaveOTP(ctx context.Context, email string, otpCode string, ttl time.Duration) error {
	data := otpData{
		Code:     otpCode,
		Attempts: 0,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, r.otpKey(email), jsonData, ttl).Err()
}

func (r *authRepository) GetOTP(ctx context.Context, email string) (string, int, error) {
	jsonData, err := r.rdb.Get(ctx, r.otpKey(email)).Bytes()
	if err != nil {
		return "", 0, err
	}
	var data otpData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return "", 0, err
	}
	return data.Code, data.Attempts, nil
}

func (r *authRepository) IncrementOTPAttempts(ctx context.Context, email string) (int, error) {
	key := r.otpKey(email)
	jsonData, err := r.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return 0, err
	}
	var data otpData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return 0, err
	}
	data.Attempts++

	ttl, err := r.rdb.TTL(ctx, key).Result()
	if err != nil || ttl <= 0 {
		return data.Attempts, r.rdb.Del(ctx, key).Err()
	}

	newJson, err := json.Marshal(data)
	if err != nil {
		return data.Attempts, err
	}

	err = r.rdb.Set(ctx, key, newJson, ttl).Err()
	return data.Attempts, err
}

func (r *authRepository) DeleteOTP(ctx context.Context, email string) error {
	return r.rdb.Del(ctx, r.otpKey(email)).Err()
}

func (r *authRepository) SetOTPCooldown(ctx context.Context, email string, ttl time.Duration) error {
	return r.rdb.Set(ctx, r.cooldownKey(email), "1", ttl).Err()
}

func (r *authRepository) CheckOTPCooldown(ctx context.Context, email string) (bool, time.Duration, error) {
	key := r.cooldownKey(email)
	exists, err := r.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, 0, err
	}
	if exists == 0 {
		return false, 0, nil
	}
	ttl, err := r.rdb.TTL(ctx, key).Result()
	if err != nil {
		return true, 0, err
	}
	return true, ttl, nil
}

func (r *authRepository) SaveResetToken(ctx context.Context, token string, email string, ttl time.Duration) error {
	return r.rdb.Set(ctx, r.resetTokenKey(token), email, ttl).Err()
}

func (r *authRepository) GetResetToken(ctx context.Context, token string) (string, error) {
	return r.rdb.Get(ctx, r.resetTokenKey(token)).Result()
}

func (r *authRepository) DeleteResetToken(ctx context.Context, token string) error {
	return r.rdb.Del(ctx, r.resetTokenKey(token)).Err()
}
