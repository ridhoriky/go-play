package order

import (
	"context"

	"ne-project/src/internal/models/dto"
)

type OrderServiceItf interface {
	Checkout(ctx context.Context, buyerID string, req *dto.CheckoutRequest) ([]dto.OrderDetailResponse, error)
	GetOrderDetail(ctx context.Context, buyerID string, orderID string) (*dto.OrderDetailResponse, error)
	GetBuyerOrders(ctx context.Context, buyerID string, params *dto.GetOrdersQuery) (*dto.OrderListResponse, error)
	GetSellerOrders(ctx context.Context, storeID string, params *dto.GetOrdersQuery) (*dto.OrderListResponse, error)
	GetSellerOrderDetail(ctx context.Context, storeID string, orderID string) (*dto.OrderDetailResponse, error)
	CancelOrder(ctx context.Context, buyerID string, orderID string) error
	ConfirmReceived(ctx context.Context, buyerID string, orderID string) error
	UpdateOrderStatus(ctx context.Context, storeID string, orderID string, status string) error
}
