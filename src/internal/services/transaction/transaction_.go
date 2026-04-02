package transaction

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
)

func (s *transactionService) Checkout(
	ctx context.Context, req *dto.CreateTransactionRequest,
) (*dto.TransactionDetailResponse, error) {

	result, err := s.transactionRepository.Checkout(ctx, req)
	if err != nil {
		return nil, err
	}

	return buildTransactionDetailResponse(result), nil
}

func buildTransactionDetailResponse(result *entity.TransactionWithDetails) *dto.TransactionDetailResponse {
	items := make([]dto.TransactionItemResponse, 0, len(result.Items))

	for _, item := range result.Items {
		items = append(items, dto.TransactionItemResponse{
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Quantity:    item.Quantity,
			Price:       item.Price.InexactFloat64(),
			Subtotal:    item.Subtotal.InexactFloat64(),
		})
	}

	return &dto.TransactionDetailResponse{
		ID:          result.ID,
		TotalAmount: result.TotalAmount.InexactFloat64(),
		Status:      string(result.Status),
		CreatedAt:   result.CreatedAt,
		Items:       items,
	}
}

func (s *transactionService) GetTransactionByID(ctx context.Context, id string) (*dto.TransactionDetailResponse, error) {
	result, err := s.transactionRepository.GetTransactionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return buildTransactionDetailResponse(result), nil
}
