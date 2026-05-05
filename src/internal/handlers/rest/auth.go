package rest

import (
	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService auth.AuthServiceItf
}

func NewAuthHandler(authService auth.AuthServiceItf) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password, returns access and refresh tokens
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

	loginResp, err := h.authService.Login(ctx, req.Email, req.Password, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		helpers.ResponseError(c, err)
		return
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
// @Success      200 {object} dto.RegisterResponse
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
// @Description  Generate new access token using refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshTokenRequest true "Refresh token"
// @Success      200 {object} dto.AuthTokenResponse
// @Failure      400 {object} dto.Error
// @Failure      401 {object} dto.Error
// @Failure      500 {object} dto.Error
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	tokenResp, err := h.authService.RefreshToken(ctx, req.RefreshToken, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Token refreshed successfully", tokenResp)
}

// Logout godoc
// @Summary      User logout
// @Description  Invalidate refresh token and logout user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LogoutRequest true "Refresh token"
// @Success      200 {object} dto.MessageResponse
// @Failure      400 {object} dto.Error
// @Failure      401 {object} dto.Error
// @Failure      500 {object} dto.Error
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		helpers.ResponseError(c, dto.NewError(http.StatusUnauthorized, preference.ErrMissingAuthHeader))
		return
	}

	if err := h.authService.Logout(ctx, userID.(string), req.RefreshToken); err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Logout successful", nil)
}
