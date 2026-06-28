package rest

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/cart"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService cart.CartServiceItf
}

func NewCartHandler(cartService cart.CartServiceItf) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// AddToCart godoc
// @Summary Add item to cart
// @Tags Cart
// @Accept json
// @Produce json
// @Param request body dto.AddToCartRequest true "Cart item details"
// @Success 201 {object} dto.CartItemResponse
// @Router /cart [post]
func (h *CartHandler) AddToCart(c *gin.Context) {
	ctx := c.Request.Context()
	buyerID := c.GetString("user_id")

	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	res, err := h.cartService.AddToCart(ctx, buyerID, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Item added to cart", res)
}

// GetCart godoc
// @Summary Get my cart
// @Tags Cart
// @Produce json
// @Success 200 {object} dto.CartResponse
// @Router /cart [get]
func (h *CartHandler) GetCart(c *gin.Context) {
	ctx := c.Request.Context()
	buyerID := c.GetString("user_id")

	res, err := h.cartService.GetCart(ctx, buyerID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", res)
}

// UpdateCartItem godoc
// @Summary Update cart item quantity
// @Tags Cart
// @Accept json
// @Produce json
// @Param id path string true "Cart Item ID"
// @Param request body dto.UpdateCartRequest true "Update details"
// @Success 200 {string} string "Success"
// @Router /cart/{id} [put]
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	ctx := c.Request.Context()
	buyerID := c.GetString("user_id")
	cartID := c.Param("id")

	var req dto.UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	err := h.cartService.UpdateQuantity(ctx, buyerID, cartID, &req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Cart item updated successfully", nil)
}

// RemoveFromCart godoc
// @Summary Remove item from cart
// @Tags Cart
// @Produce json
// @Param id path string true "Cart Item ID"
// @Success 200 {string} string "Success"
// @Router /cart/{id} [delete]
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	ctx := c.Request.Context()
	buyerID := c.GetString("user_id")
	cartID := c.Param("id")

	err := h.cartService.RemoveFromCart(ctx, buyerID, cartID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Item removed from cart", nil)
}

// ClearCart godoc
// @Summary Clear all items from cart
// @Tags Cart
// @Produce json
// @Success 200 {string} string "Success"
// @Router /cart [delete]
func (h *CartHandler) ClearCart(c *gin.Context) {
	ctx := c.Request.Context()
	buyerID := c.GetString("user_id")

	err := h.cartService.ClearCart(ctx, buyerID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Cart cleared successfully", nil)
}
