package payment

import (
	"context"

	"ne-project/src/internal/models/dto"
)

type PaymentServiceItf interface {
	CreatePayment(ctx context.Context, orderID string, method string) (*dto.PaymentResult, error)
	HandleCallback(ctx context.Context, ref string, status string) error
	GetPaymentStatus(ctx context.Context, orderID string) (*dto.PaymentStatus, error)
}
