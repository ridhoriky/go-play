package report

import (
	"context"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/utils/validation"

	"github.com/jmoiron/sqlx"
)

type ReportRepositoryItf interface {
	GetSummary(ctx context.Context, dr validation.DateRange) (dto.SummaryResponse, error)
}
type reportRepository struct {
	db *sqlx.DB
}

func NewReportRepository(db *sqlx.DB) ReportRepositoryItf {
	return &reportRepository{
		db: db,
	}
}
