package auth

// import (
// 	"context"
// 	"ne-project/src/internal/models/dto"
// 	"ne-project/src/internal/preference"
// 	"net/http"
// )

// func (s *auth) RequestChallenge(ctx context.Context, req *dto.RequestChallengeRequest) (*dto.ChallengeResponse, error) {
// 	email := req.Email
// 	publicKey := req.PublicKey

// 	if !email || !publicKey {
// 		return &dto.Error{
// 			Code:    http.StatusBadRequest,
// 			Message: preference.ErrProductStockNegative,
// 		}
// 	}

// 	if !validatePublicKey(publicKey) {
// 		return &dto.Error{
// 			Code:    http.StatusBadRequest,
// 			Message: "Invalid public key",
// 		}
// 	}

// 	existing := s.userRepository.GetByEmail(ctx, email)
// 	if existing {
// 		return &dto.Error{
// 			Code:    http.StatusConflict,
// 			Message: "Email already exists",
// 		}
// 	}
// 	return nil, nil
// }

// func (s *auth) VerifyChallenge(ctx context.Context, req *dto.VerifyChallengeRequest) (*dto.ChallengeResponse, error) {
// 	return nil, nil
// }

// func (s *auth) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
// 	return nil, nil
// }

// func (s *auth) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
// 	return nil, nil
// }

// func (s *auth) Logout(ctx context.Context, req *dto.LogoutRequest) (*dto.LogoutResponse, error) {
// 	return nil, nil
// }

// func (s *auth) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
// 	return nil, nil
// }
