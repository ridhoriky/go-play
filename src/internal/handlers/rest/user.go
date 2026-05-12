package rest

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService user.UserServiceItf
}

func NewUserHandler(userService user.UserServiceItf) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetAllUser godoc
// @Summary Get all users
// @Description Get a list of all users with pagination
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {object} dto.UserListResponse
// @Failure 400 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /users [get]
func (h *UserHandler) GetAllUser(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.GetUsersQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidQueryParams))
		return
	}

	users, err := h.userService.GetAllUsers(ctx, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", users)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get a single user by their unique ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} entity.User
// @Failure 400 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidUserID))
		return
	}
	user, err := h.userService.GetUserByID(ctx, id.String())
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", user)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User information"
// @Success 201 {object} entity.User
// @Failure 400 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()

	var user dto.CreateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	createdUser, err := h.userService.CreateUser(ctx, &user)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusCreated, "User created successfully", createdUser)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update user information by their unique ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dto.UpdateUserRequest true "Updated user information"
// @Success 200 {object} entity.User
// @Failure 400 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidUserID))
		return
	}

	var user dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	updatedUser, err := h.userService.UpdateUser(ctx, id.String(), &user)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "User updated successfully", updatedUser)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by their unique ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.Error
// @Failure 404 {object} dto.Error
// @Failure 500 {object} dto.Error
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidUserID))
		return
	}
	if err := h.userService.DeleteUser(ctx, id.String()); err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "User deleted successfully", nil)
}
