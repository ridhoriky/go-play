package token

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"ne-project/src/internal/models/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type TokenOptions struct {
	SecretAccessToken   string        `yaml:"secret_token" env:"JWT_SECRET_ACCESS_TOKEN" env-default:"secret_access_token"`
	SecretRefreshToken  string        `yaml:"secret_refresh_token" env:"JWT_SECRET_REFRESH_TOKEN" env-default:"secret_refresh_token"`
	ExpiredToken        time.Duration `yaml:"expired_token" env:"JWT_EXPIRED_ACCESS_TOKEN" env-default:"15m"`
	ExpiredRefreshToken time.Duration `yaml:"expired_refresh_token" env:"JWT_EXPIRED_REFRESH_TOKEN" env-default:"168h"`
}

type TokenDetails struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int64  `json:"expiresAt"`
	ExpiresRt    int64  `json:"expiresRt"`
}

type AccessTokenClaims struct {
	Authorized bool   `json:"authorized"`
	UserID     string `json:"user_id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	AccessUUID string `json:"access_uuid"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	TokenType   string `json:"token_type"`
	RefreshUUID string `json:"refresh_uuid"`
	UserID      string `json:"user_id"`
	jwt.RegisteredClaims
}

type Token struct {
	log                 zerolog.Logger
	secretAccessToken   []byte
	secretRefreshToken  []byte
	expiredToken        time.Duration
	expiredRefreshToken time.Duration
}

func InitToken(log zerolog.Logger, opt TokenOptions) *Token {
	if len(strings.TrimSpace(opt.SecretAccessToken)) == 0 || len(strings.TrimSpace(opt.SecretRefreshToken)) == 0 {
		log.Panic().Msgf("Environment variable jwt access or refresh key is not set")
	}
	return &Token{
		log:                 log,
		secretAccessToken:   []byte(strings.TrimSpace(opt.SecretAccessToken)),
		secretRefreshToken:  []byte(strings.TrimSpace(opt.SecretRefreshToken)),
		expiredToken:        opt.ExpiredToken,
		expiredRefreshToken: opt.ExpiredRefreshToken,
	}

}

func (a *Token) CreateTokens(user *entity.User) (*TokenDetails, error) {
	now := time.Now()
	td := &TokenDetails{
		ExpiresAt: time.Now().Add(a.expiredToken).Unix(),
		ExpiresRt: time.Now().Add(a.expiredRefreshToken).Unix(),
	}

	accessUUID := uuid.NewString()
	refreshUUID := uuid.NewString()

	accessClaims := AccessTokenClaims{
		Authorized: true,
		UserID:     user.ID,
		Name:       user.Name,
		Role:       user.Role,
		AccessUUID: accessUUID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.ExpiresAt, 0)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	refreshClaims := RefreshTokenClaims{
		TokenType:   "refresh",
		RefreshUUID: refreshUUID,
		UserID:      user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.ExpiresRt, 0)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString(a.secretAccessToken)
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString(a.secretRefreshToken)
	if err != nil {
		return nil, err
	}

	td.AccessToken = accessString
	td.RefreshToken = refreshString

	return td, nil
}

func (a *Token) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.secretAccessToken, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}

	return claims, nil
}

func (a *Token) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.secretRefreshToken, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("invalid refresh token")
	}

	return claims, nil
}

func (a *Token) HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func (a *Token) ExtractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if strings.TrimSpace(authHeader) == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.TrimSpace(parts[0]) != "Bearer" {
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	return strings.TrimSpace(parts[1]), nil
}

func (a *Token) ValidateToken(r *http.Request) error {
	tokenString, err := a.ExtractTokenFromHeader(r)
	if err != nil {
		return err
	}

	_, err = a.ValidateAccessToken(tokenString)
	return err
}

func (a *Token) ValidateRefreshTokenFromRequest(r *http.Request, token string) error {
	if token == "" {
		return fmt.Errorf("refresh token is required")
	}
	_, err := a.ValidateRefreshToken(token)
	return err
}
