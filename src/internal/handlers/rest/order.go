package rest

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/order"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService order.OrderServiceItf
}

func NewOrderHandler(orderService order.OrderServiceItf) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// Checkout godoc
// @Summary Checkout cart items
// @Tags Order (Buyer)
// @Accept json
// @Produce json
// @Param request body dto.CheckoutRequest true "Checkout details"
// @Success 201 {array} dto.OrderDetailResponse
// @Router /orders [post]
func (h *OrderHandler) Checkout(c *gin.Context) {
	ctx := c.Request.Context()
	buyerID := c.GetString("user_id")

	var req dto.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	res, err := h.orderService.Checkout(ctx, buyerID, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Checkout successful", res)
}

// GetBuyerOrders godoc
// @Summary Get buyer orders
// @Tags Order (Buyer)
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page limit"
// @Param status query string false "Order status"
// @Success 200 {object} dto.OrderListResponse
// @Router /orders [get]
func (h *OrderHandler) GetBuyerOrders(c *gin.Context) {
	ctx := c.Request.Context()
	buyerID := c.GetString("user_id")

	var query dto.GetOrdersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		query = dto.GetOrdersQuery{Page: 1, Limit: 10, SortBy: "created_at", SortDir: "DESC"}
	}
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}

	res, err := h.orderService.GetBuyerOrders(ctx, buyerID, &query)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", res)
}

// GetBuyerOrderDetail godoc
// @Summary Get buyer order detail
// @Tags Order (Buyer)
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} dto.OrderDetailResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) GetBuyerOrderDetail(c *gin.Context) {
	ctx := c.Request.Context()
	buyerID := c.GetString("user_id")
	orderID := c.Param("id")

	res, err := h.orderService.GetOrderDetail(ctx, buyerID, orderID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", res)
}

// CancelOrder godoc
// @Summary Cancel order
// @Tags Order (Buyer)
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {string} string "Success"
// @Router /orders/{id}/cancel [patch]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	ctx := c.Request.Context()
	buyerID := c.GetString("user_id")
	orderID := c.Param("id")

	if err := h.orderService.CancelOrder(ctx, buyerID, orderID); err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Order canceled successfully", nil)
}

// ConfirmReceived godoc
// @Summary Confirm order received
// @Tags Order (Buyer)
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {string} string "Success"
// @Router /orders/{id}/confirm [patch]
func (h *OrderHandler) ConfirmReceived(c *gin.Context) {
	ctx := c.Request.Context()
	buyerID := c.GetString("user_id")
	orderID := c.Param("id")

	if err := h.orderService.ConfirmReceived(ctx, buyerID, orderID); err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Order confirmed successfully", nil)
}

// GetSellerOrders godoc
// @Summary Get seller orders
// @Tags Order (Seller)
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page limit"
// @Param status query string false "Order status"
// @Success 200 {object} dto.OrderListResponse
// @Router /seller/orders [get]
func (h *OrderHandler) GetSellerOrders(c *gin.Context) {
	ctx := c.Request.Context()
	storeID := c.GetString("store_id")

	var query dto.GetOrdersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		query = dto.GetOrdersQuery{Page: 1, Limit: 10, SortBy: "created_at", SortDir: "DESC"}
	}
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}

	res, err := h.orderService.GetSellerOrders(ctx, storeID, &query)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", res)
}

// GetSellerOrderDetail godoc
// @Summary Get seller order detail
// @Tags Order (Seller)
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} dto.OrderDetailResponse
// @Router /seller/orders/{id} [get]
func (h *OrderHandler) GetSellerOrderDetail(c *gin.Context) {
	ctx := c.Request.Context()
	storeID := c.GetString("store_id")
	orderID := c.Param("id")

	res, err := h.orderService.GetSellerOrderDetail(ctx, storeID, orderID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", res)
}

// UpdateSellerOrderStatus godoc
// @Summary Update order status
// @Tags Order (Seller)
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body dto.UpdateOrderStatusRequest true "Status"
// @Success 200 {string} string "Success"
// @Router /seller/orders/{id}/status [patch]
func (h *OrderHandler) UpdateSellerOrderStatus(c *gin.Context) {
	ctx := c.Request.Context()
	storeID := c.GetString("store_id")
	orderID := c.Param("id")

	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	if err := h.orderService.UpdateOrderStatus(ctx, storeID, orderID, req.Status); err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Order status updated successfully", nil)
}
