package token

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"ne-project/src/internal/models/entity"

	"github.com/golang-jwt/jwt/v5"
)

type TokenOptions struct {
	AccessPrivateKeyPEM  string        `yaml:"access_private_key_pem" env:"JWT_ACCESS_PRIVATE_KEY"`
	AccessPublicKeyPEM   string        `yaml:"access_public_key_pem" env:"JWT_ACCESS_PUBLIC_KEY"`
	RefreshPrivateKeyPEM string        `yaml:"refresh_private_key_pem" env:"JWT_REFRESH_PRIVATE_KEY"`
	RefreshPublicKeyPEM  string        `yaml:"refresh_public_key_pem" env:"JWT_REFRESH_PUBLIC_KEY"`
	ExpiredToken         time.Duration `yaml:"expired_token" env:"JWT_EXPIRED_ACCESS_TOKEN" env-default:"15m"`
	ExpiredRefreshToken  time.Duration `yaml:"expired_refresh_token" env:"JWT_EXPIRED_REFRESH_TOKEN" env-default:"168h"`
	CookieDomain         string        `yaml:"cookie_domain" env:"JWT_COOKIE_DOMAIN" env-default:""`
	CookieSecure         bool          `yaml:"cookie_secure" env:"JWT_COOKIE_SECURE" env-default:"true"`
}

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresRt    int64  `json:"expires_rt"`
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

type TokenServiceItf interface {
	CreateTokens(user *entity.User) (*TokenDetails, error)
	ValidateAccessToken(tokenString string) (*AccessTokenClaims, error)
	ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error)
	HashToken(token string) string
}

type tokenService struct {
	accessPrivateKey    *ecdsa.PrivateKey
	accessPublicKey     *ecdsa.PublicKey
	refreshPrivateKey   *ecdsa.PrivateKey
	refreshPublicKey    *ecdsa.PublicKey
	expiredToken        time.Duration
	expiredRefreshToken time.Duration
}

func NewTokenService(opt *TokenOptions) (TokenServiceItf, error) {
	accessPriv, accessPub, err := parseECKeyPair(opt.AccessPrivateKeyPEM, opt.AccessPublicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("access token keys: %w", err)
	}

	refreshPriv, refreshPub, err := parseECKeyPair(opt.RefreshPrivateKeyPEM, opt.RefreshPublicKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("refresh token keys: %w", err)
	}

	return &tokenService{
		accessPrivateKey:    accessPriv,
		accessPublicKey:     accessPub,
		refreshPrivateKey:   refreshPriv,
		refreshPublicKey:    refreshPub,
		expiredToken:        opt.ExpiredToken,
		expiredRefreshToken: opt.ExpiredRefreshToken,
	}, nil
}
