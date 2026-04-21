package auth

import (
	"ne-project/src/internal/services/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandlerItf interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	RefreshToken(c *gin.Context)
	Logout(c *gin.Context)
}

type authHandler struct {
	authService auth.AuthServiceItf
}

func NewAuthHandler(authService auth.AuthServiceItf) AuthHandlerItf {
	return &authHandler{
		authService: authService,
	}
}
