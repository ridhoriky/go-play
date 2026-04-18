package user

import "github.com/gin-gonic/gin"

type UserHandlerItf interface {
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	CreateUserBulk(c *gin.Context)
}

// type userHandler struct {
// 	userService user.UserServiceItf
// }

// func NewUserHandler(userService user.UserServiceItf) UserHandlerItf {
// 	return &userHandler{
// 		userService: userService,
// 	}
// }
