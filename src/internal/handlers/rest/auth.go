package rest

import (
	"net/http"
	"time"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService  auth.AuthServiceItf
	cookieDomain string
	cookieSecure bool
}

//nolint:gosec // These are cookie settings, not hardcoded credentials
const (
	refreshTokenCookieName = "refresh_token"
	refreshTokenCookiePath = "/api/v1/auth"
)

func NewAuthHandler(authService auth.AuthServiceItf, cookieDomain string, cookieSecure bool) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		cookieDomain: cookieDomain,
		cookieSecure: cookieSecure,
	}
}

func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, refreshToken string, expiresRt int64) {
	maxAge := max(int(time.Until(time.Unix(expiresRt, 0)).Seconds()), 0)
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		refreshTokenCookieName,
		refreshToken,
		maxAge,
		refreshTokenCookiePath,
		h.cookieDomain,
		h.cookieSecure,
		true, // HttpOnly
	)
}

func (h *AuthHandler) clearRefreshTokenCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(refreshTokenCookieName, "", -1, refreshTokenCookiePath, h.cookieDomain, h.cookieSecure, true)
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login credentials"
// @Success      200 {object} dto.LoginResponse
// @Failure      400 {object} dto.Error
// @Failure      401 {object} dto.Error
// @Failure      500 {object} dto.Error
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	result, err := h.authService.Login(ctx, req.Email, req.Password, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	h.setRefreshTokenCookie(c, result.RefreshToken, result.ExpiresRt)

	loginResp := dto.LoginResponse{
		AuthTokenResponse: dto.AuthTokenResponse{
			AccessToken: result.AccessToken,
			ExpiresAt:   result.ExpiresAt,
			User:        result.User,
		},
		Message: "Login successful",
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Login successful", loginResp)
}

// Register godoc
// @Summary      Register new user
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "User registration data"
// @Success      201 {object} dto.RegisterResponse
// @Failure      400 {object} dto.Error
// @Failure      409 {object} dto.Error
// @Failure      500 {object} dto.Error
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	registerResp, err := h.authService.Register(ctx, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "Registration successful", registerResp)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Generate new access token using the refresh token stored in HttpOnly cookie.
// @Tags         auth
// @Produce      json
// @Success      200 {object} dto.AuthTokenResponse
// @Failure      401 {object} dto.Error
// @Failure      500 {object} dto.Error
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	refreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil || refreshToken == "" {
		helpers.ResponseError(c, dto.NewError(http.StatusUnauthorized, preference.ErrInvalidRefreshToken))
		return
	}

	result, err := h.authService.RefreshToken(ctx, refreshToken, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		h.clearRefreshTokenCookie(c)
		helpers.ResponseError(c, err)
		return
	}

	h.setRefreshTokenCookie(c, result.RefreshToken, result.ExpiresRt)

	tokenResp := dto.AuthTokenResponse{
		AccessToken: result.AccessToken,
		ExpiresAt:   result.ExpiresAt,
		User:        result.User,
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Token refreshed successfully", tokenResp)
}

// Logout godoc
// @Summary      User logout
// @Description  Revoke refresh token and clear the HttpOnly cookie
// @Tags         auth
// @Produce      json
// @Success      200 {object} dto.APIResponse
// @Failure      401 {object} dto.Error
// @Failure      500 {object} dto.Error
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		helpers.ResponseError(c, dto.NewError(http.StatusUnauthorized, preference.ErrMissingAuthHeader))
		return
	}

	refreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil || refreshToken == "" {
		h.clearRefreshTokenCookie(c)
		helpers.ResponseSuccess(c, http.StatusOK, "Logout successful", nil)
		return
	}

	if err := h.authService.Logout(ctx, userID.(string), refreshToken); err != nil {
		helpers.ResponseError(c, err)
		return
	}

	h.clearRefreshTokenCookie(c)
	helpers.ResponseSuccess(c, http.StatusOK, "Logout successful", nil)
}

// VerifyEmail godoc
// @Summary      Verify email with OTP
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.VerifyEmailRequest true "Email and OTP details"
// @Success      200 {object} dto.APIResponse
// @Failure      400 {object} dto.Error
// @Router       /auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	if err := h.authService.VerifyEmail(ctx, req.Email, req.OTPCode); err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Email successfully verified.", nil)
}

// ResendOTP godoc
// @Summary      Resend OTP code
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.ResendOTPRequest true "Email details"
// @Success      200 {object} dto.APIResponse
// @Failure      400 {object} dto.Error
// @Router       /auth/resend-otp [post]
func (h *AuthHandler) ResendOTP(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.ResendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	if err := h.authService.ResendOTP(ctx, req.Email); err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Verification code resent successfully.", nil)
}

// GoogleLogin godoc
// @Summary      Google OAuth Login
// @Description  Authenticate using Google ID Token. Returns access token in body; refresh token is set as HttpOnly cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.GoogleLoginRequest true "Google ID token"
// @Success      200 {object} dto.LoginResponse
// @Failure      400 {object} dto.Error
// @Router       /auth/google [post]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.GoogleLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	result, err := h.authService.LoginGoogle(ctx, req.IDToken, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	h.setRefreshTokenCookie(c, result.RefreshToken, result.ExpiresRt)

	loginResp := dto.LoginResponse{
		AuthTokenResponse: dto.AuthTokenResponse{
			AccessToken: result.AccessToken,
			ExpiresAt:   result.ExpiresAt,
			User:        result.User,
		},
		Message: "Google login successful",
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Google login successful", loginResp)
}

// ForgotPassword godoc
// @Summary      Request password reset link
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.ForgotPasswordRequest true "User email"
// @Success      200 {object} dto.APIResponse
// @Router       /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	if err := h.authService.ForgotPassword(ctx, req.Email); err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "If the email is registered, a password reset link has been sent.", nil)
}

// ResetPassword godoc
// @Summary      Reset password using token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.ResetPasswordRequest true "Token and new password"
// @Success      200 {object} dto.APIResponse
// @Failure      400 {object} dto.Error
// @Router       /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	if err := h.authService.ResetPassword(ctx, req.Token, req.NewPassword); err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Password successfully updated. You can now log in.", nil)
}
