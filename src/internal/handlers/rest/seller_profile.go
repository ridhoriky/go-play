package rest

import (
	"context"
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/store"
	"ne-project/src/internal/services/user"

	"github.com/gin-gonic/gin"
)

type SellerProfileHandler struct {
	userService  user.UserServiceItf
	storeService store.StoreServiceItf
}

func NewSellerProfileHandler(userService user.UserServiceItf, storeService store.StoreServiceItf) *SellerProfileHandler {
	return &SellerProfileHandler{
		userService:  userService,
		storeService: storeService,
	}
}

// GetSellerProfile godoc
// @Summary      Get seller profile
// @Description  Get seller profile including user and store info
// @Tags         seller-profile
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.SellerProfileResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /seller/profile [get]
func (h *SellerProfileHandler) GetSellerProfile(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	userEntity, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	storeRes, err := h.storeService.GetMyStore(ctx, userID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	userRes := dto.UserResponse{
		ID:        userEntity.ID,
		Name:      userEntity.Name,
		Email:     userEntity.Email,
		Role:      userEntity.Role,
		IsActive:  userEntity.IsActive,
		AvatarURL: userEntity.AvatarURL,
		Phone:     userEntity.Phone,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}

	res := dto.SellerProfileResponse{
		User:  userRes,
		Store: *storeRes,
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Success", res)
}

// UpdateSellerProfile godoc
// @Summary      Update seller profile
// @Description  Update seller profile including user and store info
// @Tags         seller-profile
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.UpdateSellerProfileRequest true "Update profile request"
// @Success      200  {object}  dto.SellerProfileResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /seller/profile [put]
func (h *SellerProfileHandler) UpdateSellerProfile(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	var req dto.UpdateSellerProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	if err := h.updateUserFromProfileReq(ctx, userID, &req); err != nil {
		helpers.ResponseError(c, err)
		return
	}

	storeRes, err := h.updateStoreFromProfileReq(ctx, userID, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	userEntity, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	userRes := dto.UserResponse{
		ID:        userEntity.ID,
		Name:      userEntity.Name,
		Email:     userEntity.Email,
		Role:      userEntity.Role,
		IsActive:  userEntity.IsActive,
		AvatarURL: userEntity.AvatarURL,
		Phone:     userEntity.Phone,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}

	res := dto.SellerProfileResponse{
		User:  userRes,
		Store: *storeRes,
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Profile updated successfully", res)
}

func (h *SellerProfileHandler) updateUserFromProfileReq(ctx context.Context, userID string, req *dto.UpdateSellerProfileRequest) error {
	updateUserReq := &dto.UpdateUserRequest{}
	var hasUserUpdate bool
	if req.Name != "" {
		updateUserReq.Name = &req.Name
		hasUserUpdate = true
	}
	if req.Phone != "" {
		updateUserReq.Phone = &req.Phone
		hasUserUpdate = true
	}
	if req.AvatarURL != "" {
		updateUserReq.AvatarURL = &req.AvatarURL
		hasUserUpdate = true
	}

	if hasUserUpdate {
		_, err := h.userService.UpdateUser(ctx, userID, updateUserReq)
		return err
	}
	return nil
}

func (h *SellerProfileHandler) updateStoreFromProfileReq(ctx context.Context, userID string, req *dto.UpdateSellerProfileRequest) (*dto.StoreResponse, error) {
	updateStoreReq := &dto.UpdateStoreRequest{}
	var hasStoreUpdate bool
	if req.StoreName != "" {
		updateStoreReq.StoreName = req.StoreName
		hasStoreUpdate = true
	}
	if req.Description != "" {
		updateStoreReq.Description = req.Description
		hasStoreUpdate = true
	}
	if req.LogoURL != "" {
		updateStoreReq.LogoURL = req.LogoURL
		hasStoreUpdate = true
	}
	if req.BannerURL != "" {
		updateStoreReq.BannerURL = req.BannerURL
		hasStoreUpdate = true
	}

	if !hasStoreUpdate {
		return h.storeService.GetMyStore(ctx, userID)
	}

	existingStore, err := h.storeService.GetMyStore(ctx, userID)
	if err != nil {
		return nil, err
	}
	if updateStoreReq.StoreName == "" {
		updateStoreReq.StoreName = existingStore.StoreName
	}
	if updateStoreReq.Description == "" {
		updateStoreReq.Description = existingStore.Description
	}
	if updateStoreReq.LogoURL == "" {
		updateStoreReq.LogoURL = existingStore.LogoURL
	}
	if updateStoreReq.BannerURL == "" {
		updateStoreReq.BannerURL = existingStore.BannerURL
	}

	return h.storeService.UpdateStore(ctx, userID, updateStoreReq)
}

// GetStoreStats godoc
// @Summary      Get store statistics
// @Description  Get store quick statistics
// @Tags         seller-profile
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.StoreStats
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /seller/store/stats [get]
func (h *SellerProfileHandler) GetStoreStats(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	storeRes, err := h.storeService.GetMyStore(ctx, userID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	stats, err := h.storeService.GetStoreStats(ctx, storeRes.ID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Success", stats)
}
