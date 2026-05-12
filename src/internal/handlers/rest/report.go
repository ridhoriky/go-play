package rest

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/report"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService report.ReportServiceItf
}

func NewReportHandler(reportService report.ReportServiceItf) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// GetReports godoc
// @Summary      List Report
// @Description  Get list of reports
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param		 period			query		string	false	"Period"	default(today)
// @Param		 date_from		query		string	false	"Date from"
// @Param		 date_to		query		string	false	"Date to"
// @Success      200  {object}  	dto.APIResponse{data=[]dto.SummaryResponse}
// @Failure      400  {object} 		dto.APIResponse
// @Failure      404  {object}  	dto.APIResponse
// @Router       /reports [get]
func (h *ReportHandler) GetReports(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.GetReportQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidQueryParams))
		return
	}
	reports, err := h.reportService.GetSummary(ctx, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", reports)
}

// GetTopProduct godoc
// @Summary      List Top Products
// @Description  Get list of top products
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param		 period			query		string	false	"Period"	default(today)
// @Param		 date_from		query		string	false	"Date from"
// @Param		 date_to		query		string	false	"Date to"
// @Param		 limit			query		int		false	"Limit"	default(10)
// @Success      200  {object}  	dto.APIResponse{data=[]dto.TopProductsResponse}
// @Failure      400  {object} 		dto.APIResponse
// @Failure      404  {object}  	dto.APIResponse
// @Router       /reports/top-products [get]
func (h *ReportHandler) GetTopProducts(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.GetTopProductsQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidQueryParams))
		return
	}
	reports, err := h.reportService.GetTopProducts(ctx, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", reports)
}
