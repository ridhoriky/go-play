package rest

import (
	"net/http"
	"strconv"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/report"

	"github.com/gin-gonic/gin"
)

type SellerReportHandler struct {
	sellerReportService report.SellerReportServiceItf
}

func NewSellerReportHandler(sellerReportService report.SellerReportServiceItf) *SellerReportHandler {
	return &SellerReportHandler{
		sellerReportService: sellerReportService,
	}
}

// GetSalesSummary godoc
// @Summary      Get seller sales summary
// @Description  Get seller sales summary
// @Tags         seller-reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        period  query     string  false  "Period (today, this_week, this_month)"  default(this_week)
// @Param        date_from query   string  false  "Date from (YYYY-MM-DD)"
// @Param        date_to   query   string  false  "Date to (YYYY-MM-DD)"
// @Success      200     {object}  dto.APIResponse{data=dto.SellerSalesSummary}
// @Failure      400     {object}  dto.APIResponse
// @Failure      401     {object}  dto.APIResponse
// @Router       /seller/reports/summary [get]
func (h *SellerReportHandler) GetSalesSummary(c *gin.Context) {
	ctx := c.Request.Context()
	storeID := c.GetString("store_id")

	var req dto.GetSellerReportQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidQueryParams))
		return
	}

	summary, err := h.sellerReportService.GetSalesSummary(ctx, storeID, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Success", summary)
}

// GetTopProducts godoc
// @Summary      Get seller top products
// @Description  Get seller top products
// @Tags         seller-reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        period  query     string  false  "Period (today, this_week, this_month)"  default(this_week)
// @Param        date_from query   string  false  "Date from (YYYY-MM-DD)"
// @Param        date_to   query   string  false  "Date to (YYYY-MM-DD)"
// @Param        limit   query     int     false  "Limit"  default(10)
// @Param        offset  query     int     false  "Offset" default(0)
// @Success      200     {object}  dto.APIResponse{data=[]dto.SellerTopProduct}
// @Failure      400     {object}  dto.APIResponse
// @Failure      401     {object}  dto.APIResponse
// @Router       /seller/reports/top-products [get]
func (h *SellerReportHandler) GetTopProducts(c *gin.Context) {
	ctx := c.Request.Context()
	storeID := c.GetString("store_id")

	var req dto.GetSellerReportQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidQueryParams))
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit := 10
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	offset := 0
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	products, err := h.sellerReportService.GetTopProducts(ctx, storeID, &req, limit, offset)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Success", products)
}

// GetDashboard godoc
// @Summary      Get seller dashboard
// @Description  Get seller dashboard with summary, top products, recent orders
// @Tags         seller-reports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        period  query     string  false  "Period (today, this_week, this_month)"  default(this_week)
// @Success      200     {object}  dto.APIResponse{data=dto.SellerDashboard}
// @Failure      401     {object}  dto.APIResponse
// @Router       /seller/reports/dashboard [get]
func (h *SellerReportHandler) GetDashboard(c *gin.Context) {
	ctx := c.Request.Context()
	storeID := c.GetString("store_id")

	var req dto.GetSellerReportQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidQueryParams))
		return
	}

	dashboard, err := h.sellerReportService.GetDashboard(ctx, storeID, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Success", dashboard)
}
