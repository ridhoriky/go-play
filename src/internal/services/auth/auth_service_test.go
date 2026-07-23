package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"net/http"
	"testing"
	"time"

	"ne-project/src/internal/config/token"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/repositories/auth"
	"ne-project/src/internal/repositories/user"
	"ne-project/src/internal/utils/hash"
	"ne-project/src/internal/utils/mailer"

	"github.com/rs/zerolog"
)

// ─── Struct and Type Declarations ───────────────────────────────────────────

type mockUserRepository struct {
	user.UserRepositoryItf
	GetAllFunc     func(ctx context.Context, query *dto.GetUsersQuery) ([]entity.User, int, error)
	GetByIDFunc    func(ctx context.Context, id string) (*entity.User, error)
	GetByEmailFunc func(ctx context.Context, email string) (*entity.User, error)
	CreateFunc     func(ctx context.Context, user *entity.User) error
	UpdateFunc     func(ctx context.Context, id string, user *entity.User) error
	DeleteFunc     func(ctx context.Context, id string) error
}

type mockAuthRepository struct {
	auth.AuthRepositoryItf
	SaveRefreshTokenFunc      func(ctx context.Context, tokenHash string, userID string, userAgent string, ipAddress string, expiresAt time.Time) error
	GetRefreshTokenByHashFunc func(ctx context.Context, tokenHash string) (*entity.RefreshToken, error)
	RevokeRefreshTokenFunc    func(ctx context.Context, tokenHash string) error
	RevokeAllUserTokensFunc   func(ctx context.Context, userID string) error
	SaveOTPFunc               func(ctx context.Context, email string, otpCode string, ttl time.Duration) error
	GetOTPFunc                func(ctx context.Context, email string) (string, int, error)
	IncrementOTPAttemptsFunc  func(ctx context.Context, email string) (int, error)
	DeleteOTPFunc             func(ctx context.Context, email string) error
	SetOTPCooldownFunc        func(ctx context.Context, email string, ttl time.Duration) error
	CheckOTPCooldownFunc      func(ctx context.Context, email string) (bool, time.Duration, error)
	SaveResetTokenFunc        func(ctx context.Context, token string, email string, ttl time.Duration) error
	GetResetTokenFunc         func(ctx context.Context, token string) (string, error)
	DeleteResetTokenFunc      func(ctx context.Context, token string) error
}

// ─── Mock Method Implementations ─────────────────────────────────────────────

func (m *mockUserRepository) GetAll(ctx context.Context, query *dto.GetUsersQuery) ([]entity.User, int, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, query)
	}
	return nil, 0, nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, sql.ErrNoRows
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, sql.ErrNoRows
}

func (m *mockUserRepository) Create(ctx context.Context, user *entity.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *mockUserRepository) Update(ctx context.Context, id string, user *entity.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, user)
	}
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *mockAuthRepository) SaveRefreshToken(ctx context.Context, tokenHash string, userID string, userAgent string, ipAddress string, expiresAt time.Time) error {
	if m.SaveRefreshTokenFunc != nil {
		return m.SaveRefreshTokenFunc(ctx, tokenHash, userID, userAgent, ipAddress, expiresAt)
	}
	return nil
}

func (m *mockAuthRepository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
	if m.GetRefreshTokenByHashFunc != nil {
		return m.GetRefreshTokenByHashFunc(ctx, tokenHash)
	}
	return nil, auth.ErrRefreshTokenNotFound
}

func (m *mockAuthRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	if m.RevokeRefreshTokenFunc != nil {
		return m.RevokeRefreshTokenFunc(ctx, tokenHash)
	}
	return nil
}

func (m *mockAuthRepository) RevokeAllUserTokens(ctx context.Context, userID string) error {
	if m.RevokeAllUserTokensFunc != nil {
		return m.RevokeAllUserTokensFunc(ctx, userID)
	}
	return nil
}

func (m *mockAuthRepository) SaveOTP(ctx context.Context, email string, otpCode string, ttl time.Duration) error {
	if m.SaveOTPFunc != nil {
		return m.SaveOTPFunc(ctx, email, otpCode, ttl)
	}
	return nil
}

func (m *mockAuthRepository) GetOTP(ctx context.Context, email string) (string, int, error) {
	if m.GetOTPFunc != nil {
		return m.GetOTPFunc(ctx, email)
	}
	return "", 0, nil
}

func (m *mockAuthRepository) IncrementOTPAttempts(ctx context.Context, email string) (int, error) {
	if m.IncrementOTPAttemptsFunc != nil {
		return m.IncrementOTPAttemptsFunc(ctx, email)
	}
	return 0, nil
}

func (m *mockAuthRepository) DeleteOTP(ctx context.Context, email string) error {
	if m.DeleteOTPFunc != nil {
		return m.DeleteOTPFunc(ctx, email)
	}
	return nil
}

func (m *mockAuthRepository) SetOTPCooldown(ctx context.Context, email string, ttl time.Duration) error {
	if m.SetOTPCooldownFunc != nil {
		return m.SetOTPCooldownFunc(ctx, email, ttl)
	}
	return nil
}

func (m *mockAuthRepository) CheckOTPCooldown(ctx context.Context, email string) (bool, time.Duration, error) {
	if m.CheckOTPCooldownFunc != nil {
		return m.CheckOTPCooldownFunc(ctx, email)
	}
	return false, 0, nil
}

func (m *mockAuthRepository) SaveResetToken(ctx context.Context, token string, email string, ttl time.Duration) error {
	if m.SaveResetTokenFunc != nil {
		return m.SaveResetTokenFunc(ctx, token, email, ttl)
	}
	return nil
}

func (m *mockAuthRepository) GetResetToken(ctx context.Context, token string) (string, error) {
	if m.GetResetTokenFunc != nil {
		return m.GetResetTokenFunc(ctx, token)
	}
	return "", nil
}

func (m *mockAuthRepository) DeleteResetToken(ctx context.Context, token string) error {
	if m.DeleteResetTokenFunc != nil {
		return m.DeleteResetTokenFunc(ctx, token)
	}
	return nil
}

// ─── Setup helper ───────────────────────────────────────────────────────────

// generateTestECKeyPEM generates a test ECDSA P-256 key pair and returns PEM-encoded private and public keys.
func generateTestECKeyPEM(t testing.TB) (privPEM, pubPEM string) {
	t.Helper()
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate test EC key: %v", err)
	}
	privBytes, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		t.Fatalf("failed to marshal EC private key: %v", err)
	}
	pubBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal EC public key: %v", err)
	}
	privPEM = string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}))
	pubPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}))
	return
}

func setupTestService(tb testing.TB, userRepo user.UserRepositoryItf, authRepo auth.AuthRepositoryItf) AuthServiceItf {
	tb.Helper()
	nopLogger := zerolog.Nop()

	accessPrivPEM, accessPubPEM := generateTestECKeyPEM(tb)
	refreshPrivPEM, refreshPubPEM := generateTestECKeyPEM(tb)

	tokenSvc, err := token.NewTokenService(&token.TokenOptions{
		AccessPrivateKeyPEM:  accessPrivPEM,
		AccessPublicKeyPEM:   accessPubPEM,
		RefreshPrivateKeyPEM: refreshPrivPEM,
		RefreshPublicKeyPEM:  refreshPubPEM,
		ExpiredToken:         15 * time.Minute,
		ExpiredRefreshToken:  7 * 24 * time.Hour,
	})
	if err != nil {
		tb.Fatalf("setupTestService: failed to init token service: %v", err)
	}

	mailerSvc := mailer.NewMailer(&nopLogger, "", 587, "", "", "", "")

	return NewAuthService(userRepo, authRepo, nil, tokenSvc, mailerSvc, "test-client-id", "https://test.example.com")
}

// ─── Test Cases: Login ───────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	ctx := t.Context()
	plainPassword := "Secret123!"
	hashedPwd, _ := hash.HashPassword(plainPassword)

	userRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{
				ID:         "user-1",
				Email:      "budi@example.com",
				Password:   hashedPwd,
				IsActive:   true,
				IsVerified: true,
				Role:       "user",
			}, nil
		},
	}
	authRepo := &mockAuthRepository{
		SaveRefreshTokenFunc: func(ctx context.Context, tokenHash string, userID string, userAgent string, ipAddress string, expiresAt time.Time) error {
			return nil
		},
	}

	svc := setupTestService(t, userRepo, authRepo)
	resp, err := svc.Login(ctx, "budi@example.com", plainPassword, "Chrome", "127.0.0.1")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp == nil || resp.AccessToken == "" {
		t.Error("expected access token to be generated")
	}
}

func TestLogin_ErrorUserNotFound(t *testing.T) {
	ctx := t.Context()
	plainPassword := "Secret123!"

	userRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return nil, sql.ErrNoRows
		},
	}
	authRepo := &mockAuthRepository{}

	svc := setupTestService(t, userRepo, authRepo)
	_, err := svc.Login(ctx, "nonexistent@example.com", plainPassword, "Chrome", "127.0.0.1")

	if err == nil {
		t.Error("expected error, got nil")
	}
	var appErr *dto.Error
	if errors.As(err, &appErr) && appErr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, appErr.Code)
	}
}

func TestLogin_ErrorInactiveAccount(t *testing.T) {
	ctx := t.Context()
	plainPassword := "Secret123!"
	hashedPwd, _ := hash.HashPassword(plainPassword)

	userRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{
				ID:         "user-1",
				Email:      "budi@example.com",
				Password:   hashedPwd,
				IsActive:   false,
				IsVerified: true,
			}, nil
		},
	}
	authRepo := &mockAuthRepository{}

	svc := setupTestService(t, userRepo, authRepo)
	_, err := svc.Login(ctx, "budi@example.com", plainPassword, "Chrome", "127.0.0.1")

	if err == nil {
		t.Error("expected error, got nil")
	}
	var appErr *dto.Error
	if errors.As(err, &appErr) && appErr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, appErr.Code)
	}
}

func TestLogin_ErrorUnverifiedEmail(t *testing.T) {
	ctx := t.Context()
	plainPassword := "Secret123!"
	hashedPwd, _ := hash.HashPassword(plainPassword)

	userRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{
				ID:         "user-1",
				Email:      "budi@example.com",
				Password:   hashedPwd,
				IsActive:   true,
				IsVerified: false,
			}, nil
		},
	}
	authRepo := &mockAuthRepository{}

	svc := setupTestService(t, userRepo, authRepo)
	_, err := svc.Login(ctx, "budi@example.com", plainPassword, "Chrome", "127.0.0.1")

	if err == nil {
		t.Error("expected error, got nil")
	}
	var appErr *dto.Error
	if errors.As(err, &appErr) {
		if appErr.Code != http.StatusForbidden {
			t.Errorf("expected status %d, got %d", http.StatusForbidden, appErr.Code)
		}
		if appErr.Message != "Email is not verified." {
			t.Errorf("expected error message 'Email is not verified.', got '%s'", appErr.Message)
		}
	}
}

func TestLogin_ErrorWrongPassword(t *testing.T) {
	ctx := t.Context()
	plainPassword := "Secret123!"
	hashedPwd, _ := hash.HashPassword(plainPassword)

	userRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{
				ID:         "user-1",
				Email:      "budi@example.com",
				Password:   hashedPwd,
				IsActive:   true,
				IsVerified: true,
			}, nil
		},
	}
	authRepo := &mockAuthRepository{}

	svc := setupTestService(t, userRepo, authRepo)
	_, err := svc.Login(ctx, "budi@example.com", "WrongPwd123!", "Chrome", "127.0.0.1")

	if err == nil {
		t.Error("expected error, got nil")
	}
	var appErr *dto.Error
	if errors.As(err, &appErr) && appErr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, appErr.Code)
	}
}

// ─── Test Cases: Register ───────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	ctx := t.Context()
	userRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return nil, sql.ErrNoRows
		},
		CreateFunc: func(ctx context.Context, user *entity.User) error {
			return nil
		},
	}
	otpSaved := false
	cooldownSet := false
	authRepo := &mockAuthRepository{
		SaveOTPFunc: func(ctx context.Context, email string, otpCode string, ttl time.Duration) error {
			otpSaved = true
			return nil
		},
		SetOTPCooldownFunc: func(ctx context.Context, email string, ttl time.Duration) error {
			cooldownSet = true
			return nil
		},
	}

	svc := setupTestService(t, userRepo, authRepo)
	req := &dto.RegisterRequest{
		Name:     "Budi",
		Email:    "budi@example.com",
		Password: "Password123!",
	}
	resp, err := svc.Register(ctx, req)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp == nil || resp.Message == "" {
		t.Error("expected valid registration response")
	}
	if !otpSaved || !cooldownSet {
		t.Error("expected OTP and cooldown to be stored")
	}
}

func TestRegister_ErrorEmailAlreadyRegistered(t *testing.T) {
	ctx := t.Context()
	userRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: "user-1", Email: "budi@example.com"}, nil
		},
	}
	authRepo := &mockAuthRepository{}

	svc := setupTestService(t, userRepo, authRepo)
	req := &dto.RegisterRequest{
		Name:     "Budi",
		Email:    "budi@example.com",
		Password: "Password123!",
	}
	_, err := svc.Register(ctx, req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	var appErr *dto.Error
	if errors.As(err, &appErr) && appErr.Code != http.StatusConflict {
		t.Errorf("expected conflict status, got %d", appErr.Code)
	}
}

func TestRegister_ErrorWeakPassword(t *testing.T) {
	ctx := t.Context()
	userRepo := &mockUserRepository{}
	authRepo := &mockAuthRepository{}

	svc := setupTestService(t, userRepo, authRepo)
	req := &dto.RegisterRequest{
		Name:     "Budi",
		Email:    "budi@example.com",
		Password: "123",
	}
	_, err := svc.Register(ctx, req)

	if err == nil {
		t.Error("expected error, got nil")
	}
	var appErr *dto.Error
	if errors.As(err, &appErr) && appErr.Code != http.StatusBadRequest {
		t.Errorf("expected bad request status, got %d", appErr.Code)
	}
}

// ─── Test Cases: Verify Email ────────────────────────────────────────────────

func TestVerifyEmail_Success(t *testing.T) {
	ctx := t.Context()
	otpDeleted := false
	userUpdated := false

	userRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: "user-1", Email: "budi@example.com", IsVerified: false}, nil
		},
		UpdateFunc: func(ctx context.Context, id string, user *entity.User) error {
			if user.IsVerified {
				userUpdated = true
			}
			return nil
		},
	}
	authRepo := &mockAuthRepository{
		GetOTPFunc: func(ctx context.Context, email string) (string, int, error) {
			return "123456", 0, nil
		},
		DeleteOTPFunc: func(ctx context.Context, email string) error {
			otpDeleted = true
			return nil
		},
	}

	svc := setupTestService(t, userRepo, authRepo)
	err := svc.VerifyEmail(ctx, "budi@example.com", "123456")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !userUpdated {
		t.Error("expected user to be updated to verified")
	}
	if !otpDeleted {
		t.Error("expected OTP to be deleted from Redis")
	}
}

func TestVerifyEmail_ErrorOTPExpired(t *testing.T) {
	ctx := t.Context()
	userRepo := &mockUserRepository{}
	authRepo := &mockAuthRepository{
		GetOTPFunc: func(ctx context.Context, email string) (string, int, error) {
			return "", 0, errors.New("key not found")
		},
	}

	svc := setupTestService(t, userRepo, authRepo)
	err := svc.VerifyEmail(ctx, "budi@example.com", "123456")

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestVerifyEmail_ErrorWrongOTP_RemainingAttempts(t *testing.T) {
	ctx := t.Context()
	authRepo := &mockAuthRepository{
		GetOTPFunc: func(ctx context.Context, email string) (string, int, error) {
			return "123456", 1, nil
		},
		IncrementOTPAttemptsFunc: func(ctx context.Context, email string) (int, error) {
			return 2, nil
		},
	}
	userRepo := &mockUserRepository{}

	svc := setupTestService(t, userRepo, authRepo)
	err := svc.VerifyEmail(ctx, "budi@example.com", "wrongcode")

	if err == nil {
		t.Error("expected error, got nil")
	}
	var appErr *dto.Error
	if errors.As(err, &appErr) {
		if appErr.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, appErr.Code)
		}
		if appErr.Message != "Invalid verification code. 3 attempts remaining." && appErr.Message != "Invalid or expired verification code." {
			t.Errorf("unexpected error message: %s", appErr.Message)
		}
	}
}

func TestVerifyEmail_ErrorWrongOTP_TooManyAttempts(t *testing.T) {
	ctx := t.Context()
	otpDeleted := false

	authRepo := &mockAuthRepository{
		GetOTPFunc: func(ctx context.Context, email string) (string, int, error) {
			return "123456", 4, nil
		},
		IncrementOTPAttemptsFunc: func(ctx context.Context, email string) (int, error) {
			return 5, nil
		},
		DeleteOTPFunc: func(ctx context.Context, email string) error {
			otpDeleted = true
			return nil
		},
	}
	userRepo := &mockUserRepository{}

	svc := setupTestService(t, userRepo, authRepo)
	err := svc.VerifyEmail(ctx, "budi@example.com", "wrongcode")

	if err == nil {
		t.Error("expected error, got nil")
	}
	if !otpDeleted {
		t.Error("expected OTP to be deleted after 5 attempts")
	}
	var appErr *dto.Error
	if errors.As(err, &appErr) && appErr.Message != "Too many failed attempts. Verification code has been invalidated." {
		t.Errorf("unexpected message: %s", appErr.Message)
	}
}

// ─── Test Cases: Resend OTP ──────────────────────────────────────────────────

func TestResendOTP_Success(t *testing.T) {
	ctx := t.Context()
	userRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: "user-1", Email: "budi@example.com", IsVerified: false}, nil
		},
	}
	otpSaved := false
	cooldownSet := false
	authRepo := &mockAuthRepository{
		CheckOTPCooldownFunc: func(ctx context.Context, email string) (bool, time.Duration, error) {
			return false, 0, nil
		},
		SaveOTPFunc: func(ctx context.Context, email string, otpCode string, ttl time.Duration) error {
			otpSaved = true
			return nil
		},
		SetOTPCooldownFunc: func(ctx context.Context, email string, ttl time.Duration) error {
			cooldownSet = true
			return nil
		},
	}

	svc := setupTestService(t, userRepo, authRepo)
	err := svc.ResendOTP(ctx, "budi@example.com")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !otpSaved || !cooldownSet {
		t.Error("expected new OTP to be saved and cooldown set")
	}
}

func TestResendOTP_ErrorAlreadyVerified(t *testing.T) {
	ctx := t.Context()
	userRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: "user-1", Email: "budi@example.com", IsVerified: true}, nil
		},
	}
	authRepo := &mockAuthRepository{}

	svc := setupTestService(t, userRepo, authRepo)
	err := svc.ResendOTP(ctx, "budi@example.com")

	if err == nil {
		t.Error("expected error, got nil")
	}
	var appErr *dto.Error
	if errors.As(err, &appErr) && appErr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, appErr.Code)
	}
}

func TestResendOTP_ErrorCooldownActive(t *testing.T) {
	ctx := t.Context()
	userRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: "user-1", Email: "budi@example.com", IsVerified: false}, nil
		},
	}
	authRepo := &mockAuthRepository{
		CheckOTPCooldownFunc: func(ctx context.Context, email string) (bool, time.Duration, error) {
			return true, 45 * time.Second, nil
		},
	}

	svc := setupTestService(t, userRepo, authRepo)
	err := svc.ResendOTP(ctx, "budi@example.com")

	if err == nil {
		t.Error("expected error, got nil")
	}
	var appErr *dto.Error
	if errors.As(err, &appErr) && appErr.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429 status, got %d", appErr.Code)
	}
}

// ─── Test Cases: Forgot Password ─────────────────────────────────────────────

func TestForgotPassword_Success(t *testing.T) {
	ctx := t.Context()
	userRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: "user-1", Email: "budi@example.com"}, nil
		},
	}
	resetTokenSaved := false
	authRepo := &mockAuthRepository{
		SaveResetTokenFunc: func(ctx context.Context, token string, email string, ttl time.Duration) error {
			resetTokenSaved = true
			return nil
		},
	}

	svc := setupTestService(t, userRepo, authRepo)
	err := svc.ForgotPassword(ctx, "budi@example.com")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !resetTokenSaved {
		t.Error("expected reset token to be saved")
	}
}

func TestForgotPassword_SuccessSimulationNonExistent(t *testing.T) {
	ctx := t.Context()
	userRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return nil, sql.ErrNoRows
		},
	}
	authRepo := &mockAuthRepository{}

	svc := setupTestService(t, userRepo, authRepo)
	err := svc.ForgotPassword(ctx, "unknown@example.com")

	if err != nil {
		t.Errorf("expected no error (simulation), got %v", err)
	}
}

// ─── Test Cases: Reset Password ──────────────────────────────────────────────

func TestResetPassword_Success(t *testing.T) {
	ctx := t.Context()
	userRepo := &mockUserRepository{
		GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: "user-1", Email: "budi@example.com"}, nil
		},
		UpdateFunc: func(ctx context.Context, id string, user *entity.User) error {
			return nil
		},
	}
	tokenDeleted := false
	tokensRevoked := false
	authRepo := &mockAuthRepository{
		GetResetTokenFunc: func(ctx context.Context, token string) (string, error) {
			return "budi@example.com", nil
		},
		DeleteResetTokenFunc: func(ctx context.Context, token string) error {
			tokenDeleted = true
			return nil
		},
		RevokeAllUserTokensFunc: func(ctx context.Context, userID string) error {
			tokensRevoked = true
			return nil
		},
	}

	svc := setupTestService(t, userRepo, authRepo)
	err := svc.ResetPassword(ctx, "my-reset-token", "NewPassword123!")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !tokenDeleted {
		t.Error("expected reset token to be deleted")
	}
	if !tokensRevoked {
		t.Error("expected user refresh tokens to be revoked")
	}
}

func TestResetPassword_ErrorInvalidToken(t *testing.T) {
	ctx := t.Context()
	userRepo := &mockUserRepository{}
	authRepo := &mockAuthRepository{
		GetResetTokenFunc: func(ctx context.Context, token string) (string, error) {
			return "", errors.New("token expired")
		},
	}

	svc := setupTestService(t, userRepo, authRepo)
	err := svc.ResetPassword(ctx, "invalid-token", "NewPassword123!")

	if err == nil {
		t.Error("expected error, got nil")
	}
}

// ─── Test Cases: Refresh Token ──────────────────────────────────────────────

func TestRefreshToken_Success(t *testing.T) {
	ctx := t.Context()
	plainPassword := "Secret123!"
	hashedPwd, _ := hash.HashPassword(plainPassword)
	testUser := &entity.User{
		ID:         "user-1",
		Email:      "budi@example.com",
		Password:   hashedPwd,
		IsActive:   true,
		IsVerified: true,
		Role:       "user",
	}

	accessPrivPEM, accessPubPEM := generateTestECKeyPEM(t)
	refreshPrivPEM, refreshPubPEM := generateTestECKeyPEM(t)
	tokenSvc, _ := token.NewTokenService(&token.TokenOptions{
		AccessPrivateKeyPEM:  accessPrivPEM,
		AccessPublicKeyPEM:   accessPubPEM,
		RefreshPrivateKeyPEM: refreshPrivPEM,
		RefreshPublicKeyPEM:  refreshPubPEM,
		ExpiredToken:         15 * time.Minute,
		ExpiredRefreshToken:  7 * 24 * time.Hour,
	})

	tokens, _ := tokenSvc.CreateTokens(testUser)
	hashedToken := tokenSvc.HashToken(tokens.RefreshToken)

	userRepo := &mockUserRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*entity.User, error) {
			return testUser, nil
		},
	}

	authRepo := &mockAuthRepository{
		GetRefreshTokenByHashFunc: func(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
			if tokenHash == hashedToken {
				return &entity.RefreshToken{
					ID:        tokenHash,
					UserID:    testUser.ID,
					TokenHash: tokenHash,
					ExpiresAt: time.Now().Add(1 * time.Hour),
					CreatedAt: time.Now(),
					RevokedAt: nil,
				}, nil
			}
			return nil, auth.ErrRefreshTokenNotFound
		},
		RevokeRefreshTokenFunc: func(ctx context.Context, tokenHash string) error {
			return nil
		},
		SaveRefreshTokenFunc: func(ctx context.Context, tokenHash string, userID string, userAgent string, ipAddress string, expiresAt time.Time) error {
			return nil
		},
	}

	nopLogger := zerolog.Nop()
	mailerSvc := mailer.NewMailer(&nopLogger, "", 587, "", "", "", "")
	svc := NewAuthService(userRepo, authRepo, nil, tokenSvc, mailerSvc, "test-client-id", "https://test.example.com")

	resp, err := svc.RefreshToken(ctx, tokens.RefreshToken, "Chrome", "127.0.0.1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if resp == nil || resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Error("expected new access and refresh tokens")
	}
}

func TestRefreshToken_ErrorExpired(t *testing.T) {
	ctx := t.Context()
	testUser := &entity.User{
		ID:         "user-1",
		Email:      "budi@example.com",
		IsActive:   true,
		IsVerified: true,
		Role:       "user",
	}

	accessPrivPEM, accessPubPEM := generateTestECKeyPEM(t)
	refreshPrivPEM, refreshPubPEM := generateTestECKeyPEM(t)
	tokenSvc, _ := token.NewTokenService(&token.TokenOptions{
		AccessPrivateKeyPEM:  accessPrivPEM,
		AccessPublicKeyPEM:   accessPubPEM,
		RefreshPrivateKeyPEM: refreshPrivPEM,
		RefreshPublicKeyPEM:  refreshPubPEM,
		ExpiredToken:         15 * time.Minute,
		ExpiredRefreshToken:  -1 * time.Second,
	})

	tokens, _ := tokenSvc.CreateTokens(testUser)

	userRepo := &mockUserRepository{}
	authRepo := &mockAuthRepository{}

	nopLogger := zerolog.Nop()
	mailerSvc := mailer.NewMailer(&nopLogger, "", 587, "", "", "", "")
	svc := NewAuthService(userRepo, authRepo, nil, tokenSvc, mailerSvc, "test-client-id", "https://test.example.com")

	_, err := svc.RefreshToken(ctx, tokens.RefreshToken, "Chrome", "127.0.0.1")
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestRefreshToken_ErrorReuse(t *testing.T) {
	ctx := t.Context()
	testUser := &entity.User{
		ID:         "user-1",
		Email:      "budi@example.com",
		IsActive:   true,
		IsVerified: true,
		Role:       "user",
	}

	accessPrivPEM, accessPubPEM := generateTestECKeyPEM(t)
	refreshPrivPEM, refreshPubPEM := generateTestECKeyPEM(t)
	tokenSvc, _ := token.NewTokenService(&token.TokenOptions{
		AccessPrivateKeyPEM:  accessPrivPEM,
		AccessPublicKeyPEM:   accessPubPEM,
		RefreshPrivateKeyPEM: refreshPrivPEM,
		RefreshPublicKeyPEM:  refreshPubPEM,
		ExpiredToken:         15 * time.Minute,
		ExpiredRefreshToken:  7 * 24 * time.Hour,
	})

	tokens, _ := tokenSvc.CreateTokens(testUser)
	hashedToken := tokenSvc.HashToken(tokens.RefreshToken)

	userRepo := &mockUserRepository{}

	revokeAllCalled := false
	revokedTime := time.Now().Add(-10 * time.Minute)
	authRepo := &mockAuthRepository{
		GetRefreshTokenByHashFunc: func(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
			if tokenHash == hashedToken {
				return &entity.RefreshToken{
					ID:        tokenHash,
					UserID:    testUser.ID,
					TokenHash: tokenHash,
					ExpiresAt: time.Now().Add(1 * time.Hour),
					CreatedAt: time.Now().Add(-1 * time.Hour),
					RevokedAt: &revokedTime,
				}, nil
			}
			return nil, auth.ErrRefreshTokenNotFound
		},
		RevokeAllUserTokensFunc: func(ctx context.Context, userID string) error {
			if userID == testUser.ID {
				revokeAllCalled = true
			}
			return nil
		},
	}

	nopLogger := zerolog.Nop()
	mailerSvc := mailer.NewMailer(&nopLogger, "", 587, "", "", "", "")
	svc := NewAuthService(userRepo, authRepo, nil, tokenSvc, mailerSvc, "test-client-id", "https://test.example.com")

	_, err := svc.RefreshToken(ctx, tokens.RefreshToken, "Chrome", "127.0.0.1")
	if err == nil {
		t.Error("expected error due to reuse, got nil")
	}
	if !revokeAllCalled {
		t.Error("expected RevokeAllUserTokens to be called upon reuse detection")
	}
}

func TestRefreshToken_ErrorInactiveAccount(t *testing.T) {
	ctx := t.Context()
	testUser := &entity.User{
		ID:         "user-1",
		Email:      "budi@example.com",
		IsActive:   false, // disabled
		IsVerified: true,
		Role:       "user",
	}

	accessPrivPEM, accessPubPEM := generateTestECKeyPEM(t)
	refreshPrivPEM, refreshPubPEM := generateTestECKeyPEM(t)
	tokenSvc, _ := token.NewTokenService(&token.TokenOptions{
		AccessPrivateKeyPEM:  accessPrivPEM,
		AccessPublicKeyPEM:   accessPubPEM,
		RefreshPrivateKeyPEM: refreshPrivPEM,
		RefreshPublicKeyPEM:  refreshPubPEM,
		ExpiredToken:         15 * time.Minute,
		ExpiredRefreshToken:  7 * 24 * time.Hour,
	})

	tokens, _ := tokenSvc.CreateTokens(testUser)

	userRepo := &mockUserRepository{
		GetByIDFunc: func(ctx context.Context, id string) (*entity.User, error) {
			return testUser, nil
		},
	}

	authRepo := &mockAuthRepository{
		GetRefreshTokenByHashFunc: func(ctx context.Context, tokenHash string) (*entity.RefreshToken, error) {
			return &entity.RefreshToken{
				ID:        tokenHash,
				UserID:    testUser.ID,
				TokenHash: tokenHash,
				ExpiresAt: time.Now().Add(1 * time.Hour),
				CreatedAt: time.Now(),
				RevokedAt: nil,
			}, nil
		},
		RevokeRefreshTokenFunc: func(ctx context.Context, tokenHash string) error {
			return nil
		},
	}

	nopLogger := zerolog.Nop()
	mailerSvc := mailer.NewMailer(&nopLogger, "", 587, "", "", "", "")
	svc := NewAuthService(userRepo, authRepo, nil, tokenSvc, mailerSvc, "test-client-id", "https://test.example.com")

	_, err := svc.RefreshToken(ctx, tokens.RefreshToken, "Chrome", "127.0.0.1")
	if err == nil {
		t.Error("expected error for inactive user, got nil")
	}
}

// ─── Test Cases: Logout ──────────────────────────────────────────────────────

func TestLogout_Success(t *testing.T) {
	ctx := t.Context()
	testUser := &entity.User{
		ID:         "user-1",
		Email:      "budi@example.com",
		IsActive:   true,
		IsVerified: true,
		Role:       "user",
	}

	accessPrivPEM, accessPubPEM := generateTestECKeyPEM(t)
	refreshPrivPEM, refreshPubPEM := generateTestECKeyPEM(t)
	tokenSvc, _ := token.NewTokenService(&token.TokenOptions{
		AccessPrivateKeyPEM:  accessPrivPEM,
		AccessPublicKeyPEM:   accessPubPEM,
		RefreshPrivateKeyPEM: refreshPrivPEM,
		RefreshPublicKeyPEM:  refreshPubPEM,
		ExpiredToken:         15 * time.Minute,
		ExpiredRefreshToken:  7 * 24 * time.Hour,
	})

	tokens, _ := tokenSvc.CreateTokens(testUser)
	hashedToken := tokenSvc.HashToken(tokens.RefreshToken)

	userRepo := &mockUserRepository{}
	revokeCalled := false
	authRepo := &mockAuthRepository{
		RevokeRefreshTokenFunc: func(ctx context.Context, tokenHash string) error {
			if tokenHash == hashedToken {
				revokeCalled = true
			}
			return nil
		},
	}

	nopLogger := zerolog.Nop()
	mailerSvc := mailer.NewMailer(&nopLogger, "", 587, "", "", "", "")
	svc := NewAuthService(userRepo, authRepo, nil, tokenSvc, mailerSvc, "test-client-id", "https://test.example.com")

	err := svc.Logout(ctx, testUser.ID, tokens.RefreshToken)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !revokeCalled {
		t.Error("expected RevokeRefreshToken to be called during logout")
	}
}

func TestLogout_ErrorUserMismatch(t *testing.T) {
	ctx := t.Context()
	testUser := &entity.User{
		ID:         "user-1",
		Email:      "budi@example.com",
		IsActive:   true,
		IsVerified: true,
		Role:       "user",
	}

	accessPrivPEM, accessPubPEM := generateTestECKeyPEM(t)
	refreshPrivPEM, refreshPubPEM := generateTestECKeyPEM(t)
	tokenSvc, _ := token.NewTokenService(&token.TokenOptions{
		AccessPrivateKeyPEM:  accessPrivPEM,
		AccessPublicKeyPEM:   accessPubPEM,
		RefreshPrivateKeyPEM: refreshPrivPEM,
		RefreshPublicKeyPEM:  refreshPubPEM,
		ExpiredToken:         15 * time.Minute,
		ExpiredRefreshToken:  7 * 24 * time.Hour,
	})

	tokens, _ := tokenSvc.CreateTokens(testUser)

	userRepo := &mockUserRepository{}
	authRepo := &mockAuthRepository{}

	nopLogger := zerolog.Nop()
	mailerSvc := mailer.NewMailer(&nopLogger, "", 587, "", "", "", "")
	svc := NewAuthService(userRepo, authRepo, nil, tokenSvc, mailerSvc, "test-client-id", "https://test.example.com")

	err := svc.Logout(ctx, "user-2", tokens.RefreshToken)
	if err == nil {
		t.Error("expected error due to user mismatch, got nil")
	}
}
