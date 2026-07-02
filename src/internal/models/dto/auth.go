package dto

// ─── Requests ──────────────────────────────────────────────────────────

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name"     binding:"required,min=1,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"     binding:"omitempty,oneof=admin user"`
}

type VerifyEmailRequest struct {
	Email   string `json:"email"   binding:"required,email"`
	OTPCode string `json:"otpCode" binding:"required,len=6"`
}

type ResendOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type GoogleLoginRequest struct {
	IDToken string `json:"idToken" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"       binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// ─── Responses ────────────────────────────────────────────────────────

type AuthTokenResponse struct {
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
	ExpiresAt    int64         `json:"expiresAt"`
	ExpiresRt    int64         `json:"expiresRt"`
	User         *UserResponse `json:"user"`
}

type LoginResponse struct {
	AuthTokenResponse
	Message string `json:"message"`
}

type RegisterResponse struct {
	User    *UserResponse `json:"user"`
	Message string        `json:"message"`
}
