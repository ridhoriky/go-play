package report

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/repositories/report"
)

type ReportServiceItf interface {
	GetSummary(ctx context.Context, req *dto.GetReportQuery) (*dto.SummaryResponse, error)
}

type reportService struct {
	reportRepository report.ReportRepositoryItf
}

func NewReportService(reportRepository report.ReportRepositoryItf) ReportServiceItf {
	return &reportService{
		reportRepository: reportRepository,
	}
}
