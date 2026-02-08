package dto

import (
	"errors"
	"fmt"
	"ne-project/internal/models"
	"time"
)
type TransactionResponse struct {
	ID          int                     `json:"id"`
	TotalAmount int                     `json:"total_amount"`
	CreatedAt   time.Time               `json:"created_at"`
	Details     []TransactionDetailDTO  `json:"details"`
}


type TransactionDetailDTO struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}



func (r *CheckoutRequest) Validate() error {

	if len(r.Items) == 0 {
		return errors.New("checkout items cannot be empty")
	}

	seen := map[int]struct{}{}

	for i, item := range r.Items {

		if item.ProductID <= 0 {
			return fmt.Errorf("invalid product_id at item[%d]", i)
		}

		if item.Quantity <= 0 {
			return fmt.Errorf("quantity must be greater than 0 at item[%d]", i)
		}

		if _, ok := seen[item.ProductID]; ok {
			return fmt.Errorf("duplicate product_id: %d", item.ProductID)
		}

		seen[item.ProductID] = struct{}{}
	}

	return nil
}
func ToTransactionResponse(m *models.Transaction) *TransactionResponse {

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


func ToTransactionResponses(models []models.Transaction) []TransactionResponse {

	res := make([]TransactionResponse, 0, len(models))

	for i := range models {
		res = append(res, *ToTransactionResponse(&models[i]))
	}

	return res
}