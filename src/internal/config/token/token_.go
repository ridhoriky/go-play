package token

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"ne-project/src/internal/models/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func parseECKeyPair(privPEM, pubPEM string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privPEM = strings.TrimSpace(privPEM)
	pubPEM = strings.TrimSpace(pubPEM)

	if privPEM == "" || pubPEM == "" {
		return nil, nil, errors.New("private and public key PEM must both be provided")
	}

	privBlock, _ := pem.Decode([]byte(privPEM))
	if privBlock == nil {
		return nil, nil, errors.New("failed to decode private key PEM block")
	}

	var privateKey *ecdsa.PrivateKey
	var err error
	switch privBlock.Type {
	case "EC PRIVATE KEY":
		privateKey, err = x509.ParseECPrivateKey(privBlock.Bytes)
	case "PRIVATE KEY":
		key, e := x509.ParsePKCS8PrivateKey(privBlock.Bytes)
		if e != nil {
			return nil, nil, fmt.Errorf("failed to parse PKCS8 private key: %w", e)
		}
		var ok bool
		privateKey, ok = key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, nil, errors.New("PKCS8 key is not an ECDSA key")
		}
	default:
		return nil, nil, fmt.Errorf("unexpected PEM block type for private key: %s", privBlock.Type)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse EC private key: %w", err)
	}

	pubBlock, _ := pem.Decode([]byte(pubPEM))
	if pubBlock == nil {
		return nil, nil, errors.New("failed to decode public key PEM block")
	}
	pubKey, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	ecPub, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, errors.New("public key is not an ECDSA key")
	}

	return privateKey, ecPub, nil
}

func (s *tokenService) CreateTokens(user *entity.User) (*TokenDetails, error) {
	now := time.Now()
	td := &TokenDetails{
		ExpiresAt: now.Add(s.expiredToken).Unix(),
		ExpiresRt: now.Add(s.expiredRefreshToken).Unix(),
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
			Issuer:    "greenmart-api",
			Audience:  jwt.ClaimStrings{"greenmart-client"},
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.ExpiresAt, 0)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        accessUUID,
		},
	}

	refreshClaims := RefreshTokenClaims{
		TokenType:   "refresh",
		RefreshUUID: refreshUUID,
		UserID:      user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			Issuer:    "greenmart-api",
			Audience:  jwt.ClaimStrings{"greenmart-refresh"},
			ExpiresAt: jwt.NewNumericDate(time.Unix(td.ExpiresRt, 0)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        refreshUUID,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodES256, accessClaims)
	accessString, err := accessToken.SignedString(s.accessPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims)
	refreshString, err := refreshToken.SignedString(s.refreshPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	td.AccessToken = accessString
	td.RefreshToken = refreshString

	return td, nil
}

func (s *tokenService) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.accessPublicKey, nil
	}, jwt.WithIssuedAt(), jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

func (s *tokenService) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.refreshPublicKey, nil
	}, jwt.WithIssuedAt(), jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid refresh token: wrong token type")
	}

	return claims, nil
}

func (s *tokenService) HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

func (s *tokenService) ExtractTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if strings.TrimSpace(authHeader) == "" {
		return "", errors.New("authorization header missing")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.TrimSpace(parts[0]) != "Bearer" {
		return "", errors.New("authorization header format must be Bearer {token}")
	}

	return strings.TrimSpace(parts[1]), nil
}
