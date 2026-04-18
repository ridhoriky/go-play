package report

import (
	"ne-project/src/internal/services/report"

	"github.com/gin-gonic/gin"
)

type ReportHandlerItf interface {
	GetReports(c *gin.Context)
	GetTopProducts(c *gin.Context)
}

type reportHandler struct {
	reportService report.ReportServiceItf
}

func NewReportHandler(reportService report.ReportServiceItf) ReportHandlerItf {
	return &reportHandler{
		reportService: reportService,
	}
}
