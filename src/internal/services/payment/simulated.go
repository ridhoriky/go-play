package payment

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/repositories/order"

	"github.com/rs/zerolog"
)

type simulatedPayment struct {
	orderRepo order.OrderRepositoryItf
}

func NewSimulatedPayment(orderRepo order.OrderRepositoryItf) PaymentServiceItf {
	return &simulatedPayment{
		orderRepo: orderRepo,
	}
}

func generatePaymentRef() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "PAY-" + hex.EncodeToString(b)
}

func (s *simulatedPayment) CreatePayment(ctx context.Context, orderID string, method string) (*dto.PaymentResult, error) {
	ref := generatePaymentRef()
	url := "http://localhost:8080/api/v1/payments/callback?ref=" + ref

	err := s.orderRepo.UpdatePayment(ctx, orderID, method, ref)
	if err != nil {
		return nil, err
	}

	// Auto-succeed after 5 seconds via a background goroutine
	bgCtx := context.WithoutCancel(ctx) // Use request-independent context for goroutine
	go func(paymentRef string) {
		time.Sleep(5 * time.Second)
		err := s.HandleCallback(bgCtx, paymentRef, "success")
		if err != nil {
			zerolog.Ctx(bgCtx).Error().Err(err).Str("paymentRef", paymentRef).Msg("Failed to auto-succeed simulated payment")
		} else {
			zerolog.Ctx(bgCtx).Info().Str("paymentRef", paymentRef).Msg("Simulated payment auto-succeeded")
		}
	}(ref)

	return &dto.PaymentResult{
		OrderID:    orderID,
		PaymentRef: ref,
		PaymentURL: url,
		Status:     "pending",
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}, nil
}

func (s *simulatedPayment) HandleCallback(ctx context.Context, ref string, status string) error {
	o, _, err := s.orderRepo.GetByPaymentRef(ctx, ref)
	if err != nil {
		return err
	}

	switch status {
	case "success":
		if o.Status == "pending" {
			err = s.orderRepo.UpdateStatus(ctx, o.ID, "paid")
			if err != nil {
				return err
			}
		}
	case "failed":
		if o.Status == "pending" {
			err = s.orderRepo.UpdateStatus(ctx, o.ID, "canceled")
			if err != nil {
				return err
			}
		}
	default:
		return dto.NewError(http.StatusBadRequest, "Invalid payment status callback")
	}

	return nil
}

func (s *simulatedPayment) autoSyncPaid(ctx context.Context, orderID string, o *entity.Order) {
	if err := s.orderRepo.UpdateStatus(ctx, orderID, "paid"); err != nil {
		return
	}
	updatedOrder, _, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		o.Status = "paid"
		o.UpdatedAt = time.Now()
		return
	}
	*o = *updatedOrder
}

func (s *simulatedPayment) GetPaymentStatus(ctx context.Context, orderID string) (*dto.PaymentStatus, error) {
	o, _, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Auto-sync for local dev / simulated payment:
	// If order status is pending and has a payment ref, update status to paid during polling
	if o.Status == "pending" && o.PaymentRef != nil && *o.PaymentRef != "" {
		s.autoSyncPaid(ctx, orderID, o)
	}

	var paidAt *time.Time
	if o.Status == "paid" || o.Status == "processing" || o.Status == "shipped" || o.Status == "delivered" {
		t := o.UpdatedAt
		paidAt = &t
	}

	var ref string
	if o.PaymentRef != nil {
		ref = *o.PaymentRef
	}

	return &dto.PaymentStatus{
		OrderID:       o.ID,
		Status:        o.Status,
		PaidAt:        paidAt,
		PaymentRef:    ref,
		PaymentMethod: o.PaymentMethod,
	}, nil
}
