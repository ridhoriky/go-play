package rest

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/store"

	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	storeService store.StoreServiceItf
}

func NewStoreHandler(storeService store.StoreServiceItf) *StoreHandler {
	return &StoreHandler{
		storeService: storeService,
	}
}

// CreateStore godoc
// @Summary Create a new store
// @Tags Stores
// @Accept json
// @Produce json
// @Param request body dto.CreateStoreRequest true "Store details"
// @Success 201 {object} dto.StoreResponse
// @Router /seller/stores [post]
func (h *StoreHandler) CreateStore(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	var req dto.CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	res, err := h.storeService.CreateStore(ctx, userID, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Store created successfully", res)
}

// GetMyStore godoc
// @Summary Get my store details
// @Tags Stores
// @Produce json
// @Success 200 {object} dto.StoreResponse
// @Router /seller/stores/me [get]
func (h *StoreHandler) GetMyStore(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	res, err := h.storeService.GetMyStore(ctx, userID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", res)
}

// UpdateStore godoc
// @Summary Update store details
// @Tags Stores
// @Accept json
// @Produce json
// @Param request body dto.UpdateStoreRequest true "Store details"
// @Success 200 {object} dto.StoreResponse
// @Router /seller/stores/me [put]
func (h *StoreHandler) UpdateStore(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.GetString("user_id")

	var req dto.UpdateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	res, err := h.storeService.UpdateStore(ctx, userID, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Store updated successfully", res)
}

// ListStores godoc
// @Summary List all stores
// @Tags Stores
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Param search query string false "Search query"
// @Success 200 {object} dto.StoreListResponse
// @Router /stores [get]
func (h *StoreHandler) ListStores(c *gin.Context) {
	ctx := c.Request.Context()
	var query dto.GetStoresQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidQueryParams))
		return
	}

	res, err := h.storeService.ListStores(ctx, &query)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", res)
}

// GetStoreBySlug godoc
// @Summary Get store by slug
// @Tags Stores
// @Produce json
// @Param slug path string true "Store Slug"
// @Success 200 {object} dto.StoreResponse
// @Router /stores/{slug} [get]
func (h *StoreHandler) GetStoreBySlug(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	res, err := h.storeService.GetStoreBySlug(ctx, slug)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", res)
}
