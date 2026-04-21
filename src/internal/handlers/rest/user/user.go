package user

import (
	"ne-project/src/internal/services/user"

	"github.com/gin-gonic/gin"
)

type UserHandlerItf interface {
	GetAllUser(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetUserByID(c *gin.Context)
}

type userHandler struct {
	userService user.UserServiceItf
}

func NewUserHandler(userService user.UserServiceItf) UserHandlerItf {
	return &userHandler{
		userService: userService,
	}
}
