package dto

import "time"

type PaymentResult struct {
	OrderID    string    `json:"order_id"`
	PaymentRef string    `json:"payment_ref"`
	PaymentURL string    `json:"payment_url"` // for redirect-based methods
	Status     string    `json:"status"`      // pending, success, failed
	ExpiresAt  time.Time `json:"expires_at"`
}

type PaymentStatus struct {
	OrderID       string     `json:"order_id"`
	Status        string     `json:"status"`
	PaidAt        *time.Time `json:"paid_at"`
	PaymentRef    string     `json:"payment_ref"`
	PaymentMethod *string    `json:"payment_method,omitempty"`
}

type PaymentCallbackRequest struct {
	PaymentRef string `json:"payment_ref" validate:"required"`
	Status     string `json:"status" validate:"required,oneof=success failed"`
}

type CreatePaymentRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required,oneof=bank_transfer e_wallet credit_card"`
}
