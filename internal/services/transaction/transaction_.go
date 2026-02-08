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
) ([]dto.TransactionResponse, error) {

	txs, err := s.transactionRepository.GetToday(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]dto.TransactionResponse, 0)

	for _, t := range txs {

		details := make([]dto.TransactionDetailDTO, 0)

		for _, d := range t.Details {

			details = append(details, dto.TransactionDetailDTO{
				ID: 			d.ID,
				TransactionID:	d.TransactionID,
				ProductID:   	d.ProductID,
				ProductName: 	d.ProductName,
				Quantity:    	d.Quantity,
				Subtotal:    	d.Subtotal,
			})
		}

		res = append(res, dto.TransactionResponse{
			ID:          t.ID,
			TotalAmount: t.TotalAmount,
			CreatedAt:   t.CreatedAt,
			Details:     details,
		})
	}

	return res, nil
}
