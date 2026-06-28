package rest

import (
	"net/http"
	"strconv"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/services/admin"

	"github.com/gin-gonic/gin"
)

type AdminReportHandler struct {
	service admin.AdminReportServiceItf
}

func NewAdminReportHandler(service admin.AdminReportServiceItf) *AdminReportHandler {
	return &AdminReportHandler{
		service: service,
	}
}

func (h *AdminReportHandler) GetPlatformSummary(c *gin.Context) {
	summary, err := h.service.GetPlatformSummary(c.Request.Context())
	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "platform summary retrieved successfully", summary)
}

func (h *AdminReportHandler) GetTopStores(c *gin.Context) {
	sortBy := c.Query("sort_by")
	limitStr := c.Query("limit")

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	stores, err := h.service.GetTopStores(c.Request.Context(), sortBy, limit)
	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "top stores retrieved successfully", stores)
}

func (h *AdminReportHandler) GetTopProducts(c *gin.Context) {
	sortBy := c.Query("sort_by")
	limitStr := c.Query("limit")

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	products, err := h.service.GetTopProducts(c.Request.Context(), sortBy, limit)
	if err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "top products retrieved successfully", products)
}
