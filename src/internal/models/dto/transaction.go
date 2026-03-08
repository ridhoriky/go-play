package dto

import (
	"errors"
	"fmt"
	"time"

	"ne-project/src/internal/models/entity"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionResponse struct {
	ID          string                 `json:"id"`
	TotalAmount decimal.Decimal        `json:"total_amount"`
	CreatedAt   time.Time              `json:"created_at"`
	Details     []TransactionDetailDTO `json:"details"`
}

type TransactionDetailDTO struct {
	ID            string          `json:"id"`
	TransactionID string          `json:"transaction_id"`
	ProductID     string          `json:"product_id"`
	ProductName   string          `json:"product_name"`
	Quantity      int             `json:"quantity"`
	Subtotal      decimal.Decimal `json:"subtotal"`
}

func (r *CheckoutRequest) Validate() error {

	if len(r.Items) == 0 {
		return errors.New("checkout items cannot be empty")
	}

	seen := make(map[string]struct{})

	for i, item := range r.Items {
		if item.ProductID == "" {
			return fmt.Errorf("invalid product_id at item[%d]: cannot be empty", i)
		}

		if _, err := uuid.Parse(item.ProductID); err != nil {
			return fmt.Errorf("invalid uuid format at item[%d]: %w", i, err)
		}

		if item.Quantity <= 0 {
			return fmt.Errorf("quantity must be greater than 0 at item[%d]", i)
		}

		if _, ok := seen[item.ProductID]; ok {
			return fmt.Errorf("duplicate product_id: %s", item.ProductID)
		}

		seen[item.ProductID] = struct{}{}
	}

	return nil
}
func ToTransactionResponse(m *entity.Transaction) *TransactionResponse {

	if m == nil {
		return nil
	}

	resp := &TransactionResponse{
		ID:          m.ID,
		TotalAmount: m.TotalAmount,
		CreatedAt:   m.CreatedAt,
		Details:     make([]TransactionDetailDTO, 0, len(m.Details)),
	}

	for _, d := range m.Details {
		resp.Details = append(resp.Details, TransactionDetailDTO{
			ID:            d.ID,
			TransactionID: d.TransactionID,
			ProductID:     d.ProductID,
			ProductName:   d.ProductName,
			Quantity:      d.Quantity,
			Subtotal:      d.Subtotal,
		})
	}

	return resp
}

func ToTransactionResponses(entity []entity.Transaction) []TransactionResponse {

	res := make([]TransactionResponse, 0, len(entity))

	for i := range entity {
		res = append(res, *ToTransactionResponse(&entity[i]))
	}

	return res
}
