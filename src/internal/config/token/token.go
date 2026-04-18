package token

import (
	"net/http"
	"sync"
	"time"

	_ "github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

type Token interface {
	CreateToken(r *http.Request, data any) (*TokenDetails, error)
	ValidateToken(r *http.Request) error
	ValidateRefreshToken(r *http.Request, token string) error
}

var (
	onceToken = &sync.Once{}
	tokenInst *token
)

type TokenOptions struct {
	SecretAccessToken   []byte        `yaml:"secret_access_token"`
	SecretRefreshToken  []byte        `yaml:"secret_refresh_token"`
	ExpiredToken        time.Duration `yaml:"expired_token"`
	ExpiredRefreshToken time.Duration `yaml:"expired_refresh_token"`
}

type TokenDetails struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int64  `json:"expiresAt"`
	ExpiresRt    int64  `json:"expiresRt"`
}

type token struct {
	log                 zerolog.Logger
	secretAccessToken   []byte
	secretRefreshToken  []byte
	expiredToken        time.Duration
	expiredRefreshToken time.Duration
}

func InitToken(log zerolog.Logger, opt TokenOptions) *token {
	onceToken.Do(func() {
		if len(opt.SecretAccessToken) == 0 || len(opt.SecretRefreshToken) == 0 {
			log.Panic().Msgf("Environment variable jwt access or refresh key is not set")
		}
		tokenInst = &token{
			log:                 log,
			secretAccessToken:   opt.SecretAccessToken,
			secretRefreshToken:  opt.SecretRefreshToken,
			expiredToken:        opt.ExpiredToken,
			expiredRefreshToken: opt.ExpiredRefreshToken,
		}
	})

	return tokenInst
}

// func (a *token) CreateToken(r *http.Request, data any) (*TokenDetails, error) {

// 	ctx := r.Context()
// 	td := &TokenDetails{}
// 	var err error

// 	td.ExpiresAt = time.Now().Add(a.expiredToken).Unix()

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
// 		jwt.MapClaims{
// 			"exp":         td.ExpiresAt,
// 			"access_uuid": td.AccessToken,
// 			"user_id":     user_id,
// 			"name":        username,
// 			"role":        role,
// 			"authorized":  true,
// 		})
// 	tokenString, err := token.SignedString(opt.SecretAccesToken)
// 	if err != nil {
// 		return "", nil
// 	}
// 	return tokenString, nil
// }

// func verifyToken(tokenString string) error {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return SecretAccesToken, nil
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	if !token.Valid {
// 		return fmt.Errorf("Invalid token")
// 	}
// 	return nil
// }
