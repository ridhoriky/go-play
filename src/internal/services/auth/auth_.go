package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
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
	"google.golang.org/api/idtoken"
)

func generateNumericOTP() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n.Int64())
}

func (s *authService) Login(ctx context.Context, email string, password string, userAgent string, ipAddr string) (*dto.AuthTokenResult, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil || user == nil {
		zerolog.Ctx(ctx).Warn().Str("email", email).Msg("User not found during login")
		return nil, dto.NewError(http.StatusUnauthorized, preference.ErrInvalidCredentials)
	}

	if !user.IsActive {
		zerolog.Ctx(ctx).Warn().Str("email", email).Msg("Attempt to login with disabled account")
		return nil, dto.NewError(http.StatusUnauthorized, preference.ErrInvalidCredentials)
	}

	if !user.IsVerified {
		zerolog.Ctx(ctx).Warn().Str("email", email).Msg("Attempt to login with unverified email")
		return nil, dto.NewError(http.StatusForbidden, "Email is not verified.")
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

	return &dto.AuthTokenResult{
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
	}, nil
}

func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	if req.Name == "" {
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrUserNameRequired)
	}

	if req.Email == "" {
		return nil, dto.NewError(http.StatusBadRequest, preference.ErrUserEmailRequired)
	}

	// Validate password complexity
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
		ID:         uuid.New().String(),
		Name:       req.Name,
		Email:      req.Email,
		Password:   hashedPassword,
		Role:       role,
		IsActive:   true,
		IsVerified: false,
	}

	if err := s.userRepository.Create(ctx, newUser); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to create user during registration")
		return nil, dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	// Generate OTP
	otpCode := generateNumericOTP()
	if err := s.authRepository.SaveOTP(ctx, newUser.Email, otpCode, 15*time.Minute); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to save OTP to Redis")
		return nil, dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}
	// Cooldown
	_ = s.authRepository.SetOTPCooldown(ctx, newUser.Email, 60*time.Second)

	// Send Verification Email
	s.sendVerificationEmail(ctx, newUser.Name, newUser.Email, otpCode)

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
		Message: "Registration successful. Please verify your email.",
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string, userAgent string, ipAddr string) (*dto.AuthTokenResult, error) {
	if refreshToken == "" {
		return nil, dto.NewError(http.StatusUnauthorized, preference.ErrInvalidRefreshToken)
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
	if storedToken == nil || storedToken.ExpiresAt.Before(time.Now()) {
		return nil, dto.NewError(http.StatusUnauthorized, preference.ErrInvalidRefreshToken)
	}
	if storedToken.RevokedAt != nil {
		_ = s.authRepository.RevokeAllUserTokens(ctx, storedToken.UserID)
		zerolog.Ctx(ctx).Warn().Str("user_id", storedToken.UserID).Msg("Refresh token reuse detected! All user sessions revoked.")
		return nil, dto.NewError(http.StatusUnauthorized, preference.ErrInvalidRefreshToken)
	}

	user, err := s.userRepository.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, dto.NewError(http.StatusUnauthorized, preference.ErrInvalidRefreshToken)
	}

	if !user.IsActive {
		zerolog.Ctx(ctx).Warn().Str("user_id", user.ID).Msg("Refresh token used for disabled account")
		_ = s.authRepository.RevokeRefreshToken(ctx, hashedToken)
		return nil, dto.NewError(http.StatusUnauthorized, preference.ErrInvalidCredentials)
	}

	newTokens, err := s.tokenService.CreateTokens(user)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to generate refreshed tokens")
		return nil, dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	if err = s.rotateRefreshToken(ctx, hashedToken, newTokens.RefreshToken, newTokens.ExpiresRt, user.ID, userAgent, ipAddr); err != nil {
		return nil, err
	}

	return &dto.AuthTokenResult{
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
		zerolog.Ctx(ctx).Info().Str("user_id", userID).Msg("Logout with missing refresh token — treating as success")
		return nil
	}

	claims, err := s.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Str("user_id", userID).Msg("Logout with expired/invalid refresh token — treating as success")
		return nil
	}

	if claims.UserID != userID {
		zerolog.Ctx(ctx).Warn().
			Str("user_id", userID).
			Str("token_user_id", claims.UserID).
			Msg("Logout token user mismatch")
		return dto.NewError(http.StatusUnauthorized, preference.ErrInvalidCredentials)
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
	if err := s.authRepository.RevokeRefreshToken(ctx, oldHash); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to revoke old refresh token during rotation")
		return dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}
	newHash := s.tokenService.HashToken(newRefreshToken)
	if err := s.authRepository.SaveRefreshToken(ctx, newHash, userID, userAgent, ipAddr, time.Unix(expiresRt, 0)); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to save new refresh token during rotation")
		return dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}
	return nil
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

func (s *authService) VerifyEmail(ctx context.Context, email string, otpCode string) error {
	code, attempts, err := s.authRepository.GetOTP(ctx, email)
	if err != nil {
		return dto.NewError(http.StatusBadRequest, "Invalid or expired verification code.")
	}

	if attempts >= 5 {
		_ = s.authRepository.DeleteOTP(ctx, email)
		return dto.NewError(http.StatusBadRequest, "Too many failed attempts. Verification code has been invalidated.")
	}

	if code != otpCode {
		newAttempts, incErr := s.authRepository.IncrementOTPAttempts(ctx, email)
		if incErr != nil {
			return dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
		}
		if newAttempts >= 5 {
			_ = s.authRepository.DeleteOTP(ctx, email)
			return dto.NewError(http.StatusBadRequest, "Too many failed attempts. Verification code has been invalidated.")
		}
		remaining := 5 - newAttempts
		return dto.NewError(http.StatusBadRequest, fmt.Sprintf("Invalid verification code. %d attempts remaining.", remaining))
	}

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return dto.NewError(http.StatusBadRequest, preference.ErrUserNotFound)
	}

	user.IsVerified = true
	if err := s.userRepository.Update(ctx, user.ID, user); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to update user verification status")
		return dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	_ = s.authRepository.DeleteOTP(ctx, email)
	return nil
}

func (s *authService) ResendOTP(ctx context.Context, email string) error {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.NewError(http.StatusBadRequest, "User not found.")
		}
		return dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	if user.IsVerified {
		return dto.NewError(http.StatusBadRequest, "Email is already verified.")
	}

	hasCooldown, remaining, err := s.authRepository.CheckOTPCooldown(ctx, email)
	if err != nil {
		return dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}
	if hasCooldown {
		return dto.NewError(http.StatusTooManyRequests, fmt.Sprintf("Please wait %d seconds before requesting another code.", int(remaining.Seconds())))
	}

	otpCode := generateNumericOTP()
	if err := s.authRepository.SaveOTP(ctx, email, otpCode, 15*time.Minute); err != nil {
		return dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}
	_ = s.authRepository.SetOTPCooldown(ctx, email, 60*time.Second)

	// Send Verification Email
	emailSubject := "[GreenMart] Email Verification Code"
	emailBody := fmt.Sprintf(`
		<p>Hello %s,</p>
		<p>Please use the following verification code to complete your registration:</p>
		<h2><strong>%s</strong></h2>
		<p>This code will expire in 15 minutes.</p>
	`, user.Name, otpCode)

	if err := s.mailer.SendEmail(ctx, email, emailSubject, emailBody); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to send verification email")
	}

	return nil
}

func (s *authService) LoginGoogle(ctx context.Context, idToken string, userAgent string, ipAddr string) (*dto.AuthTokenResult, error) {
	if idToken == "" {
		return nil, dto.NewError(http.StatusBadRequest, "idToken is required")
	}

	payload, err := idtoken.Validate(ctx, idToken, s.googleClientID)
	if err != nil {
		zerolog.Ctx(ctx).Warn().Err(err).Msg("Invalid Google ID Token")
		return nil, dto.NewError(http.StatusUnauthorized, "Invalid Google ID Token")
	}

	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return nil, dto.NewError(http.StatusBadRequest, "Email not found in Google ID Token")
	}

	name, ok := payload.Claims["name"].(string)
	if !ok || name == "" {
		name = "Google User"
	}

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			user = nil
		} else {
			return nil, dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
		}
	}

	if user == nil {
		user = &entity.User{
			ID:         uuid.New().String(),
			Name:       name,
			Email:      email,
			Password:   "",
			Role:       "user",
			IsActive:   true,
			IsVerified: true,
		}
		if createErr := s.userRepository.Create(ctx, user); createErr != nil {
			zerolog.Ctx(ctx).Error().Err(createErr).Msg("Failed to auto-create OAuth user")
			return nil, dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
		}
	} else if !user.IsVerified {
		user.IsVerified = true
		_ = s.userRepository.Update(ctx, user.ID, user)
	}

	tokens, err := s.tokenService.CreateTokens(user)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to generate tokens for Google login")
		return nil, dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	hashedToken := s.tokenService.HashToken(tokens.RefreshToken)
	expiresAt := time.Unix(tokens.ExpiresRt, 0)
	if err := s.authRepository.SaveRefreshToken(ctx, hashedToken, user.ID, userAgent, ipAddr, expiresAt); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to persist refresh token for Google login")
		return nil, dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	return &dto.AuthTokenResult{
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
	}, nil
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			zerolog.Ctx(ctx).Info().Str("email", email).Msg("Forgot password requested for non-existent email (simulated success)")
			return nil
		}
		return dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	token := uuid.New().String()
	if err := s.authRepository.SaveResetToken(ctx, token, user.Email, 1*time.Hour); err != nil {
		return dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	resetLink := s.frontendBaseURL + "/reset-password?token=" + token
	emailSubject := "[GreenMart] Password Reset Link"
	emailBody := fmt.Sprintf(`
		<p>Hello %s,</p>
		<p>We received a request to reset your password. Click the link below to set a new password:</p>
		<p><a href="%s" target="_blank">%s</a></p>
		<p>This link will expire in 1 hour.</p>
		<p>If you did not request this, please ignore this email.</p>
	`, user.Name, resetLink, resetLink)

	if err := s.mailer.SendEmail(ctx, user.Email, emailSubject, emailBody); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to send password reset email")
	}

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, token string, newPassword string) error {
	email, err := s.authRepository.GetResetToken(ctx, token)
	if err != nil || email == "" {
		return dto.NewError(http.StatusBadRequest, "Invalid or expired reset token.")
	}

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return dto.NewError(http.StatusBadRequest, "User not found.")
	}

	hashedPassword, err := hash.HashPassword(newPassword)
	if err != nil {
		return dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	user.Password = hashedPassword
	if err := s.userRepository.Update(ctx, user.ID, user); err != nil {
		return dto.NewError(http.StatusInternalServerError, preference.ErrInternalServer)
	}

	_ = s.authRepository.DeleteResetToken(ctx, token)
	_ = s.authRepository.RevokeAllUserTokens(ctx, user.ID)

	zerolog.Ctx(ctx).Info().Str("user_id", user.ID).Msg("Password reset successful, revoked all sessions")
	return nil
}

func (s *authService) sendVerificationEmail(ctx context.Context, name, email, otpCode string) {
	emailSubject := "[GreenMart] Email Verification Code"
	emailBody := fmt.Sprintf(`
		<p>Hello %s,</p>
		<p>Thank you for registering at GreenMart. Please use the following verification code to complete your registration:</p>
		<h2><strong>%s</strong></h2>
		<p>This code will expire in 15 minutes.</p>
	`, name, otpCode)

	if err := s.mailer.SendEmail(ctx, email, emailSubject, emailBody); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("Failed to send verification email")
	}
}
