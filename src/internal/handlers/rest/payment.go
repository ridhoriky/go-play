package rest

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/payment"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService payment.PaymentServiceItf
}

func NewPaymentHandler(paymentService payment.PaymentServiceItf) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// CreatePayment godoc
// @Summary Create payment for order
// @Tags Payment
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body dto.CreatePaymentRequest true "Payment Details"
// @Success 200 {object} dto.PaymentResult
// @Router /orders/{id}/pay [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	ctx := c.Request.Context()
	orderID := c.Param("id")
	buyerID := c.GetString("user_id") // Could add check if the order belongs to buyer in the future, if needed

	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	_ = buyerID // order ownership can be verified in a separate service layer if needed

	res, err := h.paymentService.CreatePayment(ctx, orderID, req.PaymentMethod)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Payment created successfully", res)
}

// PaymentCallback godoc
// @Summary Payment callback webhook
// @Tags Payment
// @Accept json
// @Produce json
// @Param request body dto.PaymentCallbackRequest true "Callback Details"
// @Success 200 {string} string "Success"
// @Router /payments/callback [post]
func (h *PaymentHandler) PaymentCallback(c *gin.Context) {
	ctx := c.Request.Context()

	// Normally this is POST with body, but simulated.go callback url is GET with query param ?ref=...
	// However, we can support both. Task says: POST /api/v1/payments/callback
	// But our auto-simulated URL in simulated.go is a URL, usually webhooks are POST.
	// Since our goroutine calls HandleCallback directly, this endpoint is just for manual test.

	// Check query string first (for our simulated GET-like URL, though we document it as POST)
	ref := c.Query("ref")
	status := c.Query("status")

	if ref == "" || status == "" {
		// Try to bind JSON body
		var req dto.PaymentCallbackRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
			return
		}
		ref = req.PaymentRef
		status = req.Status
	}

	if err := h.paymentService.HandleCallback(ctx, ref, status); err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Payment callback processed successfully", nil)
}

// GetPaymentStatus godoc
// @Summary Get payment status
// @Tags Payment
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} dto.PaymentStatus
// @Router /orders/{id}/payment-status [get]
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	ctx := c.Request.Context()
	orderID := c.Param("id")

	res, err := h.paymentService.GetPaymentStatus(ctx, orderID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	helpers.ResponseSuccess(c, http.StatusOK, "Success", res)
}
