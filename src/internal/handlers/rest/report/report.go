package report

import (
	"ne-project/src/internal/services/report"

	"github.com/gin-gonic/gin"
)

type ReportHandlerItf interface {
	RegisterRoutes(r *gin.RouterGroup)
}

type reportHandler struct {
	reportService report.ReportServiceItf
}

func NewTransactionHandler(reportService report.ReportServiceItf) ReportHandlerItf {
	return &reportHandler{
		reportService: reportService,
	}
}

func (h *reportHandler) RegisterRoutes(r *gin.RouterGroup) {
	reportRoutes := r.Group("/reports")
	{
		reportRoutes.GET("/summary", h.GetReports)
	}
}
