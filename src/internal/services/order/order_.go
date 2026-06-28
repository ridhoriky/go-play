package order

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/repositories/cart"
	"ne-project/src/internal/repositories/order"
	"ne-project/src/internal/repositories/product"
	"ne-project/src/internal/repositories/store"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type orderService struct {
	orderRepo   order.OrderRepositoryItf
	cartRepo    cart.CartRepositoryItf
	productRepo product.ProductRepositoryItf
	storeRepo   store.StoreRepositoryItf
}

func NewOrderService(
	orderRepo order.OrderRepositoryItf,
	cartRepo cart.CartRepositoryItf,
	productRepo product.ProductRepositoryItf,
	storeRepo store.StoreRepositoryItf,
) OrderServiceItf {
	return &orderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
		storeRepo:   storeRepo,
	}
}

func generateOrderNumber() string {
	b := make([]byte, 3) // 6 hex chars
	_, _ = rand.Read(b)
	dateStr := time.Now().Format("20060102")
	return "ORD-" + dateStr + "-" + hex.EncodeToString(b)
}

func (s *orderService) Checkout(ctx context.Context, buyerID string, req *dto.CheckoutRequest) ([]dto.OrderDetailResponse, error) {
	if len(req.CartIDs) == 0 {
		return nil, dto.NewError(http.StatusBadRequest, "cart_ids cannot be empty")
	}

	grouped, err := s.validateAndGroupCartItems(ctx, buyerID, req.CartIDs)
	if err != nil {
		return nil, err
	}

	allOrders, allOrderItems, err := s.buildOrdersAndItems(ctx, buyerID, grouped, req.ShippingAddress, req.Notes)
	if err != nil {
		return nil, err
	}

	if err := s.orderRepo.Create(ctx, allOrders, allOrderItems); err != nil {
		return nil, err
	}

	s.cleanupCartItems(ctx, req.CartIDs)

	return s.buildCheckoutResponse(ctx, allOrders, allOrderItems)
}

func (s *orderService) validateAndGroupCartItems(ctx context.Context, buyerID string, cartIDs []string) (map[string][]entity.Cart, error) {
	grouped := make(map[string][]entity.Cart)

	for _, cartID := range cartIDs {
		c, err := s.cartRepo.GetByID(ctx, cartID)
		if err != nil {
			return nil, err
		}
		if c == nil || c.BuyerID != buyerID {
			return nil, dto.NewError(http.StatusBadRequest, fmt.Sprintf("cart item %s not found or unauthorized", cartID))
		}

		prodDetail, err := s.productRepo.GetByID(ctx, c.ProductID)
		if err != nil {
			return nil, err
		}
		if prodDetail == nil {
			return nil, dto.NewError(http.StatusNotFound, fmt.Sprintf("product %s not found", c.ProductID))
		}
		if !prodDetail.IsActive {
			return nil, dto.NewError(http.StatusBadRequest, fmt.Sprintf("product %s is not active", prodDetail.Name))
		}
		if prodDetail.Stock < c.Quantity {
			return nil, dto.NewError(http.StatusBadRequest, "insufficient stock for product "+prodDetail.Name)
		}

		grouped[prodDetail.StoreID] = append(grouped[prodDetail.StoreID], *c)
	}
	return grouped, nil
}

func (s *orderService) buildOrdersAndItems(
	ctx context.Context, buyerID string, grouped map[string][]entity.Cart, shippingAddress json.RawMessage, notes string,
) ([]entity.Order, []entity.OrderItem, error) {
	var allOrders []entity.Order
	var allOrderItems []entity.OrderItem

	for storeID, cartItems := range grouped {
		orderID := uuid.New().String()
		now := time.Now()
		totalAmount := decimal.Zero
		var orderItems []entity.OrderItem

		for _, cItem := range cartItems {
			prodDetail, err := s.productRepo.GetByID(ctx, cItem.ProductID)
			if err != nil {
				return nil, nil, err
			}

			subtotal := prodDetail.Price.Mul(decimal.NewFromInt(int64(cItem.Quantity)))
			totalAmount = totalAmount.Add(subtotal)

			orderItems = append(orderItems, entity.OrderItem{
				ID:           uuid.New().String(),
				OrderID:      orderID,
				ProductID:    prodDetail.ID,
				ProductName:  prodDetail.Name,
				ProductImage: nil, // Add logic if product images exist
				Quantity:     cItem.Quantity,
				Price:        prodDetail.Price,
				Subtotal:     subtotal,
				CreatedAt:    now,
			})
		}

		o := entity.Order{
			ID:              orderID,
			BuyerID:         buyerID,
			StoreID:         storeID,
			OrderNumber:     generateOrderNumber(),
			TotalAmount:     totalAmount,
			Status:          "pending",
			ShippingAddress: shippingAddress,
			ShippingCost:    decimal.Zero,
			Notes:           nil,
			CreatedAt:       now,
			UpdatedAt:       now,
		}

		if notes != "" {
			o.Notes = &notes
		}

		allOrders = append(allOrders, o)
		allOrderItems = append(allOrderItems, orderItems...)
	}
	return allOrders, allOrderItems, nil
}

func (s *orderService) buildCheckoutResponse(ctx context.Context, allOrders []entity.Order, allOrderItems []entity.OrderItem) ([]dto.OrderDetailResponse, error) {
	var createdOrders []dto.OrderDetailResponse

	for i := range allOrders {
		o := &allOrders[i]
		st, err := s.storeRepo.GetByID(ctx, o.StoreID)
		if err != nil {
			return nil, err
		}

		var resItems []dto.OrderItemSummary
		for j := range allOrderItems {
			oi := &allOrderItems[j]
			if oi.OrderID == o.ID {
				resItems = append(resItems, dto.OrderItemSummary{
					ID:           oi.ID,
					ProductID:    oi.ProductID,
					ProductName:  oi.ProductName,
					ProductImage: oi.ProductImage,
					Quantity:     oi.Quantity,
					Price:        oi.Price,
					Subtotal:     oi.Subtotal,
				})
			}
		}

		createdOrders = append(createdOrders, dto.OrderDetailResponse{
			OrderResponse: dto.OrderResponse{
				ID:          o.ID,
				OrderNumber: o.OrderNumber,
				Store: dto.StoreSummary{
					ID:   st.ID,
					Name: st.StoreName,
					Slug: st.Slug,
				},
				Items:       resItems,
				TotalAmount: o.TotalAmount,
				Status:      o.Status,
				CreatedAt:   o.CreatedAt,
			},
			ShippingAddress: o.ShippingAddress,
			ShippingCost:    o.ShippingCost,
			PaymentMethod:   o.PaymentMethod,
			Notes:           o.Notes,
		})
	}
	return createdOrders, nil
}

func (s *orderService) cleanupCartItems(ctx context.Context, cartIDs []string) {
	for _, cartID := range cartIDs {
		_ = s.cartRepo.Delete(ctx, cartID)
	}
}

func (s *orderService) GetOrderDetail(ctx context.Context, buyerID string, orderID string) (*dto.OrderDetailResponse, error) {
	o, items, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if o.BuyerID != buyerID {
		return nil, dto.NewError(http.StatusNotFound, "order not found")
	}
	return s.mapToOrderDetailResponse(ctx, o, items)
}

func (s *orderService) GetSellerOrderDetail(ctx context.Context, storeID string, orderID string) (*dto.OrderDetailResponse, error) {
	o, items, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if o.StoreID != storeID {
		return nil, dto.NewError(http.StatusNotFound, "order not found")
	}
	return s.mapToOrderDetailResponse(ctx, o, items)
}

func (s *orderService) mapToOrderDetailResponse(ctx context.Context, o *entity.Order, items []entity.OrderItem) (*dto.OrderDetailResponse, error) {
	st, err := s.storeRepo.GetByID(ctx, o.StoreID)
	if err != nil {
		return nil, err
	}

	var resItems []dto.OrderItemSummary
	for i := range items {
		oi := &items[i]
		resItems = append(resItems, dto.OrderItemSummary{
			ID:           oi.ID,
			ProductID:    oi.ProductID,
			ProductName:  oi.ProductName,
			ProductImage: oi.ProductImage,
			Quantity:     oi.Quantity,
			Price:        oi.Price,
			Subtotal:     oi.Subtotal,
		})
	}

	return &dto.OrderDetailResponse{
		OrderResponse: dto.OrderResponse{
			ID:          o.ID,
			OrderNumber: o.OrderNumber,
			Store: dto.StoreSummary{
				ID:   st.ID,
				Name: st.StoreName,
				Slug: st.Slug,
			},
			Items:       resItems,
			TotalAmount: o.TotalAmount,
			Status:      o.Status,
			CreatedAt:   o.CreatedAt,
		},
		ShippingAddress: o.ShippingAddress,
		ShippingCost:    o.ShippingCost,
		PaymentMethod:   o.PaymentMethod,
		Notes:           o.Notes,
	}, nil
}

func (s *orderService) GetBuyerOrders(ctx context.Context, buyerID string, params *dto.GetOrdersQuery) (*dto.OrderListResponse, error) {
	orders, total, err := s.orderRepo.GetByBuyerID(ctx, buyerID, params)
	if err != nil {
		return nil, err
	}
	return s.mapOrderList(ctx, orders, total, params)
}

func (s *orderService) GetSellerOrders(ctx context.Context, storeID string, params *dto.GetOrdersQuery) (*dto.OrderListResponse, error) {
	orders, total, err := s.orderRepo.GetByStoreID(ctx, storeID, params)
	if err != nil {
		return nil, err
	}
	return s.mapOrderList(ctx, orders, total, params)
}

func (s *orderService) mapOrderList(ctx context.Context, orders []entity.Order, total int, params *dto.GetOrdersQuery) (*dto.OrderListResponse, error) {
	var res []dto.OrderResponse
	for i := range orders {
		o := &orders[i]
		st, err := s.storeRepo.GetByID(ctx, o.StoreID)
		if err != nil {
			return nil, err
		}

		_, items, _ := s.orderRepo.GetByID(ctx, o.ID) // Quick fetch items
		var resItems []dto.OrderItemSummary
		for j := range items {
			oi := &items[j]
			resItems = append(resItems, dto.OrderItemSummary{
				ID:           oi.ID,
				ProductID:    oi.ProductID,
				ProductName:  oi.ProductName,
				ProductImage: oi.ProductImage,
				Quantity:     oi.Quantity,
				Price:        oi.Price,
				Subtotal:     oi.Subtotal,
			})
		}

		res = append(res, dto.OrderResponse{
			ID:          o.ID,
			OrderNumber: o.OrderNumber,
			Store: dto.StoreSummary{
				ID:   st.ID,
				Name: st.StoreName,
				Slug: st.Slug,
			},
			Items:       resItems,
			TotalAmount: o.TotalAmount,
			Status:      o.Status,
			CreatedAt:   o.CreatedAt,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))
	if res == nil {
		res = []dto.OrderResponse{}
	}

	return &dto.OrderListResponse{
		Data: res,
		Meta: dto.PaginationMeta{
			Total:      total,
			Page:       params.Page,
			Limit:      params.Limit,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *orderService) CancelOrder(ctx context.Context, buyerID string, orderID string) error {
	o, items, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}
	if o.BuyerID != buyerID {
		return dto.NewError(http.StatusNotFound, "order not found")
	}

	if o.Status != "pending" && o.Status != "paid" {
		return dto.NewError(http.StatusBadRequest, "order cannot be canceled at this stage")
	}

	if err := s.orderRepo.CancelOrder(ctx, orderID, items); err != nil {
		return err
	}

	return nil
}

func (s *orderService) ConfirmReceived(ctx context.Context, buyerID string, orderID string) error {
	o, _, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}
	if o.BuyerID != buyerID {
		return dto.NewError(http.StatusNotFound, "order not found")
	}

	if o.Status != "shipped" {
		return dto.NewError(http.StatusBadRequest, "order can only be confirmed if shipped")
	}

	return s.orderRepo.UpdateStatus(ctx, orderID, "delivered")
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, storeID string, orderID string, status string) error {
	o, _, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}
	if o.StoreID != storeID {
		return dto.NewError(http.StatusNotFound, "order not found")
	}

	// Validate allowed transitions
	// paid -> processing -> shipped
	switch status {
	case "processing":
		if o.Status != "paid" {
			return dto.NewError(http.StatusBadRequest, "invalid status transition to processing")
		}
	case "shipped":
		if o.Status != "processing" {
			return dto.NewError(http.StatusBadRequest, "invalid status transition to shipped")
		}
	default:
		return dto.NewError(http.StatusBadRequest, "status transition not allowed from seller")
	}

	return s.orderRepo.UpdateStatus(ctx, orderID, status)
}
