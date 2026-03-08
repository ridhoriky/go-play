package transaction

import (
	"context"

	"ne-project/src/internal/models/dto"
)

func (s *transactionService) GetToday(
	ctx context.Context,
) (*dto.TodaySummaryResponse, error) {

	return s.transactionRepository.GetToday(ctx)
}
