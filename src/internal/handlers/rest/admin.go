package rest

import (
	"net/http"
	"strconv"
	"strings"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/admin"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService admin.AdminServiceItf
}

func NewAdminHandler(adminService admin.AdminServiceItf) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// ListUsers godoc
// @Summary      List all users (Admin only)
// @Description  Get a paginated list of all users, searchable by name/email, filterable by role and status
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        search query string false "Search by name or email"
// @Param        role query string false "Filter by role (admin, seller, buyer, user)"
// @Param        is_active query boolean false "Filter by active status"
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Param        sort_by query string false "Sort by (name, email, created_at)" default(created_at)
// @Param        sort_order query string false "Sort order (asc, desc)" default(desc)
// @Success      200 {object} dto.SuccessResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()

	var params dto.AdminUserListParams
	params.Search = c.Query("search")
	params.Role = c.Query("role")

	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			params.IsActive = &isActive
		}
	}

	params.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	params.SortBy = c.DefaultQuery("sort_by", "created_at")
	params.SortOrder = c.DefaultQuery("sort_order", "desc")

	users, total, err := h.adminService.ListUsers(ctx, &params)
	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	helpers.ResponseSuccessWithMeta(c, http.StatusOK, "Users retrieved successfully", users, &dto.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: (total + params.Limit - 1) / params.Limit,
	})
}

// GetUserDetail godoc
// @Summary      Get user details (Admin only)
// @Description  Get specific user details by ID
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      200 {object} dto.SuccessResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /admin/users/{id} [get]
func (h *AdminHandler) GetUserDetail(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("id")

	user, err := h.adminService.GetUserByID(ctx, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			helpers.ResponseError(c, dto.NewError(http.StatusNotFound, preference.ErrNotFound))
			return
		}
		helpers.ResponseError(c, dto.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "User details retrieved successfully", user)
}

// UpdateUser godoc
// @Summary      Update user role or status (Admin only)
// @Description  Update a user's role or active status. Cannot deactivate or change own role.
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Param        request body dto.AdminUpdateUserRequest true "Update user payload"
// @Success      200 {object} dto.SuccessResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /admin/users/{id} [put]
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	targetUserID := c.Param("id")
	currentUserID, exists := c.Get("user_id")
	if !exists {
		helpers.ResponseError(c, dto.NewError(http.StatusUnauthorized, preference.ErrMissingAuthHeader))
		return
	}

	var req dto.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	err := h.adminService.UpdateUser(ctx, currentUserID.(string), targetUserID, req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "cannot change own role") || strings.Contains(err.Error(), "cannot deactivate own account") || strings.Contains(err.Error(), "invalid role") {
			statusCode = http.StatusBadRequest
		} else if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		helpers.ResponseError(c, dto.NewError(statusCode, err.Error()))
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "User updated successfully", nil)
}

// ListSellers godoc
// @Summary      List all sellers (Admin only)
// @Description  Get a paginated list of all stores/sellers, with verification status
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        search query string false "Search by store name or owner name"
// @Param        is_verified query boolean false "Filter by verified status"
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Param        sort_by query string false "Sort by (store_name, rating_avg, total_sales, created_at)" default(created_at)
// @Param        sort_order query string false "Sort order (asc, desc)" default(desc)
// @Success      200 {object} dto.SuccessResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /admin/sellers [get]
func (h *AdminHandler) ListSellers(c *gin.Context) {
	ctx := c.Request.Context()

	var params dto.AdminSellerListParams
	params.Search = c.Query("search")

	if isVerifiedStr := c.Query("is_verified"); isVerifiedStr != "" {
		isVerified, err := strconv.ParseBool(isVerifiedStr)
		if err == nil {
			params.IsVerified = &isVerified
		}
	}

	params.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	params.SortBy = c.DefaultQuery("sort_by", "created_at")
	params.SortOrder = c.DefaultQuery("sort_order", "desc")

	sellers, total, err := h.adminService.ListSellers(ctx, &params)
	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	helpers.ResponseSuccessWithMeta(c, http.StatusOK, "Sellers retrieved successfully", sellers, &dto.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: (total + params.Limit - 1) / params.Limit,
	})
}

// VerifySeller godoc
// @Summary      Verify a seller (Admin only)
// @Description  Mark a seller store as verified
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Store ID"
// @Success      200 {object} dto.SuccessResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /admin/sellers/{id}/verify [patch]
func (h *AdminHandler) VerifySeller(c *gin.Context) {
	ctx := c.Request.Context()
	storeID := c.Param("id")

	if err := h.adminService.VerifySeller(ctx, storeID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			helpers.ResponseError(c, dto.NewError(http.StatusNotFound, "store not found"))
			return
		}
		helpers.ResponseError(c, dto.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Seller verified successfully", nil)
}

// UnverifySeller godoc
// @Summary      Unverify a seller (Admin only)
// @Description  Mark a seller store as unverified
// @Tags         admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Store ID"
// @Success      200 {object} dto.SuccessResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      403 {object} dto.ErrorResponse
// @Failure      404 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /admin/sellers/{id}/unverify [patch]
func (h *AdminHandler) UnverifySeller(c *gin.Context) {
	ctx := c.Request.Context()
	storeID := c.Param("id")

	if err := h.adminService.UnverifySeller(ctx, storeID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			helpers.ResponseError(c, dto.NewError(http.StatusNotFound, "store not found"))
			return
		}
		helpers.ResponseError(c, dto.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Seller unverified successfully", nil)
}
