package order

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
)

type OrderRepositoryItf interface {
	Create(ctx context.Context, orders []entity.Order, items []entity.OrderItem) error
	GetByID(ctx context.Context, id string) (*entity.Order, []entity.OrderItem, error)
	GetByOrderNumber(ctx context.Context, orderNumber string) (*entity.Order, []entity.OrderItem, error)
	GetByPaymentRef(ctx context.Context, paymentRef string) (*entity.Order, []entity.OrderItem, error)
	GetByBuyerID(ctx context.Context, buyerID string, params *dto.GetOrdersQuery) ([]entity.Order, int, error)
	GetByStoreID(ctx context.Context, storeID string, params *dto.GetOrdersQuery) ([]entity.Order, int, error)
	UpdateStatus(ctx context.Context, id string, status string) error
	UpdatePayment(ctx context.Context, id string, paymentMethod string, paymentRef string) error
	CancelOrder(ctx context.Context, orderID string, items []entity.OrderItem) error
	HasActiveOrdersForProduct(ctx context.Context, productID string) (bool, error)
	GetRecentOrdersByProductID(ctx context.Context, productID string, limit int) ([]dto.OrderSummary, error)
}
