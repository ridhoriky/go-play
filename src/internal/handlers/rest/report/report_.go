package report

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// GetReports godoc
// @Summary      List Report
// @Description  Get list of reports
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        search       		query   	string  false  "Filter by report name"
// @Param        category_id  		query   	string  false  "Filter by category ID (UUID)"
// @Param        min_price    		query   	number  false  "Minimum report price"
// @Param        max_price    		query   	number  false  "Maximum report price"
// @Param        in_stock     		query   	bool    false  "Filter reports that have stock"
// @Param		 page				query		int		false	"Page number"	default(1)
// @Param		 limit				query		int		false	"Page size"		default(10)
// @Param		 sort_by			query		string	false	"Sort by field"
// @Param		 sort_dir			query		string	false	"Sort direction (asc/desc)"	default(asc)
// @Success      200  {object}  	dto.APIResponse{data=[]dto.SummaryResponse}
// @Failure      400  {object} 		dto.APIResponse
// @Failure      404  {object}  	dto.APIResponse
// @Router       /reports [get]
func (h *reportHandler) GetReports(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.GetReportQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg(preference.ErrInvalidQueryParams)
		helpers.ResponseError(c.Writer, &dto.Error{
			Code:    http.StatusBadRequest,
			Message: preference.ErrInvalidQueryParams,
		})
		return
	}
	reports, err := h.reportService.GetSummary(ctx, &req)
	if err != nil {
		helpers.ResponseError(c.Writer, err)
		return
	}
	helpers.ResponseSuccess(c.Writer, http.StatusOK, "Success", reports)
}
