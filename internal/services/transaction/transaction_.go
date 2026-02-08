package transaction

import (
	"context"
	"ne-project/internal/dto"
)

func (s *transactionService) Checkout(
	ctx context.Context,
	req *dto.CheckoutRequest,
) (*dto.TransactionResponse, error) {

	txModel, err := s.transactionRepository.Checkout(ctx, req)
	if err != nil {
		return nil, err
	}

	return dto.ToTransactionResponse(txModel), nil
}

func (s *transactionService) GetToday(
	ctx context.Context,
) (*dto.TodaySummaryResponse, error) {

	
	return s.transactionRepository.GetToday(ctx)
}
