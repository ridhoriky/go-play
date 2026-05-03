package transaction

import (
	"context"
	"fmt"
	"net/http"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"

	"github.com/rs/zerolog"
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

func (s *transactionService) UpdateStatus(ctx context.Context, id string, req *dto.UpdateTransactionStatusRequest) (entity.Transaction, error) {
	trx, err := s.transactionRepository.GetTransactionByID(ctx, id)
	if err != nil {
		return entity.Transaction{}, err
	}

	newStatus := entity.TransactionStatus(req.Status)

	if err := validateStatusTransition(trx.Transaction.Status, newStatus); err != nil {
		zerolog.Ctx(ctx).Warn().Str("transaction_id", id).Msg("err validate transaction status")
		return entity.Transaction{}, err
	}

	if err := s.transactionRepository.UpdateStatus(ctx, id, newStatus, trx.Items); err != nil {
		return entity.Transaction{}, err
	}

	trx.Transaction.Status = newStatus

	return trx.Transaction, nil
}

func validateStatusTransition(current, new entity.TransactionStatus) error {
	allowed := map[entity.TransactionStatus][]entity.TransactionStatus{
		entity.TransactionStatusPending: {
			entity.TransactionStatusPaid,
			entity.TransactionStatusCancelled,
		},
		entity.TransactionStatusPaid: {
			entity.TransactionStatusCancelled,
		},
		entity.TransactionStatusCancelled: {},
	}

	for _, allowedStatus := range allowed[current] {
		if allowedStatus == new {
			return nil
		}
	}

	return dto.NewError(http.StatusUnprocessableEntity, fmt.Sprintf(
			"Transition from '%s' to '%s' is not allowed",
			current, new,
		))
}
