package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"
	"unicode"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/repositories/auth"
	"ne-project/src/internal/utils/hash"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func (s *authService) Login(ctx context.Context, email string, password string, userAgent string, ipAddr string) (*dto.LoginResponse, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil || user == nil {
		zerolog.Ctx(ctx).Warn().Str("email", email).Msg("User not found during login")
		return nil, dto.NewError(http.StatusUnauthorized, preference.ErrInvalidCredentials)
	}

	if !user.IsActive {
		zerolog.Ctx(ctx).Warn().Str("email", email).Msg("Attempt to login with disabled account")
		return nil, dto.NewError(http.StatusUnauthorized, preference.ErrInvalidCredentials)
	}

	if !hash.ComparePassword(user.Password, password) {
		zerolog.Ctx(ctx).Warn().Str("email", email).Msg("Invalid password during login")
		return nil, dto.NewError(http.StatusUnauthorized, preference.ErrInvalidCredentials)
	}

	tokens, err := s.tokenService.CreateTokens(user)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to generate tokens")
		return nil, dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	hashedToken := s.tokenService.HashToken(tokens.RefreshToken)
	expiresAt := time.Unix(tokens.ExpiresRt, 0)
	if err := s.authRepository.SaveRefreshToken(ctx, hashedToken, user.ID, userAgent, ipAddr, expiresAt); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to persist refresh token")
		return nil, dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	return &dto.LoginResponse{
		AuthTokenResponse: dto.AuthTokenResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			ExpiresAt:    tokens.ExpiresAt,
			ExpiresRt:    tokens.ExpiresRt,
			User: &dto.UserResponse{
				ID:        user.ID,
				Name:      user.Name,
				Email:     user.Email,
				Role:      user.Role,
				IsActive:  user.IsActive,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		},
		Message: "Login successful",
	}, nil
}

func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	if req.Name == "" {
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrUserNameRequired)
	}

	if req.Email == "" {
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrUserEmailRequired)
	}

	if len(req.Password) < 8 {
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrInvalidPassword)
	}

	// Check for at least one uppercase letter, one number, and one special character
	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}

	existingUser, _ := s.userRepository.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		zerolog.Ctx(ctx).Warn().Str("email", req.Email).Msg("Email already registered")
		return nil, dto.NewError(http.StatusConflict, preference.ErrEmailAlreadyRegistered)
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to hash password")
		return nil, dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	newUser := &entity.User{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     role,
		IsActive: true,
	}

	if err := s.userRepository.Create(ctx, newUser); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to create user during registration")
		return nil, dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	userResp := &dto.UserResponse{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		Role:      newUser.Role,
		IsActive:  newUser.IsActive,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}

	return &dto.RegisterResponse{
		User:    userResp,
		Message: "Registration successful",
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string, userAgent string, ipAddr string) (*dto.AuthTokenResponse, error) {
	if refreshToken == "" {
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrInvalidRefreshToken)
	}

	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Msg("Invalid refresh token")
		return nil, dto.NewError(http.StatusUnauthorized, preference.ErrInvalidRefreshToken)
	}

	hashedToken := s.tokenService.HashToken(refreshToken)
	storedToken, err := s.authRepository.GetRefreshTokenByHash(ctx, hashedToken)
	if err != nil && !errors.Is(err, auth.ErrRefreshTokenNotFound) {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to get refresh token")
		return nil, dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}
	if storedToken == nil || storedToken.RevokedAt != nil || storedToken.ExpiresAt.Before(time.Now()) {
		return nil, dto.NewError(http.StatusUnauthorized, preference.ErrInvalidRefreshToken)
	}

	user, err := s.userRepository.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	newTokens, err := s.tokenService.CreateTokens(user)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to generate refreshed tokens")
		return nil, dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	if err = s.rotateRefreshToken(ctx, hashedToken, newTokens.RefreshToken, newTokens.ExpiresRt, user.ID, userAgent, ipAddr); err != nil {
		return nil, err
	}

	return &dto.AuthTokenResponse{
		AccessToken:  newTokens.AccessToken,
		RefreshToken: newTokens.RefreshToken,
		ExpiresAt:    newTokens.ExpiresAt,
		ExpiresRt:    newTokens.ExpiresRt,
		User: &dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID string, refreshToken string) error {
	if refreshToken == "" {
		return dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Msg("Invalid refresh token on logout")
		return dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	if claims.UserID != userID {
		return dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	hashedToken := s.tokenService.HashToken(refreshToken)
	if err := s.authRepository.RevokeRefreshToken(ctx, hashedToken); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to revoke refresh token on logout")
		return dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)

	}

	zerolog.Ctx(ctx).Info().Str("user_id", userID).Msg("User logged out")
	return nil
}

func (s *authService) rotateRefreshToken(ctx context.Context, oldHash, newRefreshToken string, expiresRt int64, userID, userAgent, ipAddr string) error {
	err := s.authRepository.RevokeRefreshToken(ctx, oldHash)
	if err != nil {
		return err
	}
	newHash := s.tokenService.HashToken(newRefreshToken)
	return s.authRepository.SaveRefreshToken(ctx, newHash, userID, userAgent, ipAddr, time.Unix(expiresRt, 0))
}

func validatePassword(password string) error {
	hasUpper := false
	hasNumber := false
	hasSpecial := false
	const specialChars = "!@#$%^&*"

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasNumber = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasNumber || !hasSpecial {
		return dto.NewError(http.StatusBadRequest, "Password must contain at least one uppercase letter, one number, and one special character (!@#$%^&*)")
	}
	return nil
}
