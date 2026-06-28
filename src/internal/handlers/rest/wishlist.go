package rest

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/wishlist"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	wishlistService wishlist.WishlistServiceItf
}

func NewWishlistHandler(wishlistService wishlist.WishlistServiceItf) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
	}
}

// GetWishlist godoc
// @Summary Get my wishlist
// @Tags Wishlist
// @Produce json
// @Success 200 {object} dto.WishlistResponse
// @Router /wishlist [get]
func (h *WishlistHandler) GetWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	buyerID := c.GetString("user_id")

	var query dto.GetWishlistQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidQueryParams))
		return
	}

	res, err := h.wishlistService.GetWishlist(ctx, buyerID, &query)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success retrieve wishlist", res)
}

// AddToWishlist godoc
// @Summary Add item to wishlist
// @Tags Wishlist
// @Accept json
// @Produce json
// @Param request body dto.AddToWishlistRequest true "Wishlist product ID"
// @Success 201 {string} string "Success"
// @Router /wishlist [post]
func (h *WishlistHandler) AddToWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	buyerID := c.GetString("user_id")

	var req dto.AddToWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	err := h.wishlistService.AddToWishlist(ctx, buyerID, req.ProductID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Product added to wishlist", nil)
}

// ToggleWishlist godoc
// @Summary Toggle item in wishlist
// @Tags Wishlist
// @Accept json
// @Produce json
// @Param request body dto.AddToWishlistRequest true "Wishlist product ID"
// @Success 200 {object} map[string]interface{}
// @Router /wishlist/toggle [post]
func (h *WishlistHandler) ToggleWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	buyerID := c.GetString("user_id")

	var req dto.AddToWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	added, err := h.wishlistService.ToggleWishlist(ctx, buyerID, req.ProductID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}

	if added {
		helpers.ResponseSuccess(c, http.StatusOK, "Product added to wishlist", gin.H{"is_in_wishlist": true})
	} else {
		helpers.ResponseSuccess(c, http.StatusOK, "Product removed from wishlist", gin.H{"is_in_wishlist": false})
	}
}

// RemoveFromWishlist godoc
// @Summary Remove item from wishlist
// @Tags Wishlist
// @Produce json
// @Param productId path string true "Product ID"
// @Success 200 {string} string "Success"
// @Router /wishlist/{productId} [delete]
func (h *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	buyerID := c.GetString("user_id")
	productID := c.Param("productId")

	err := h.wishlistService.RemoveFromWishlist(ctx, buyerID, productID)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Product removed from wishlist", nil)
}
