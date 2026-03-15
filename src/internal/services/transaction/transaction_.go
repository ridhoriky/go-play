package transaction

import (
	"context"

	"ne-project/src/internal/models/dto"
)

func (s *transactionService) GetToday(
	ctx context.Context,
) (*dto.TransactionListResponse, error) {

	return s.transactionRepository.GetToday(ctx)
}
