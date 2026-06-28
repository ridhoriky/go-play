package rest

import (
	"net/http"

	"ne-project/src/internal/handlers/helpers"
	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/preference"
	"ne-project/src/internal/services/review"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService review.ReviewServiceItf
}

func NewReviewHandler(reviewService review.ReviewServiceItf) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
	}
}

// CreateReview godoc
// @Summary Create a product review
// @Tags Review
// @Accept json
// @Produce json
// @Param request body dto.CreateReviewRequest true "Review details"
// @Success 201 {object} dto.ReviewResponse
// @Router /reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	ctx := c.Request.Context()
	buyerID := c.GetString("user_id")

	var req dto.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	res, err := h.reviewService.CreateReview(ctx, buyerID, req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusCreated, "Review created successfully", res)
}

// GetProductReviews godoc
// @Summary Get product reviews
// @Tags Review
// @Produce json
// @Param id path string true "Product ID"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} dto.ReviewListResponse
// @Router /products/{id}/reviews [get]
func (h *ReviewHandler) GetProductReviews(c *gin.Context) {
	ctx := c.Request.Context()
	productID := c.Param("id")

	var query dto.GetReviewsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidQueryParams))
		return
	}

	res, err := h.reviewService.GetProductReviews(ctx, productID, query)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Success", res)
}

// ReplyToReview godoc
// @Summary Reply to a review
// @Tags Review
// @Accept json
// @Produce json
// @Param id path string true "Review ID"
// @Param request body dto.SellerReplyRequest true "Reply details"
// @Success 200 {string} string "Success"
// @Router /seller/reviews/{id}/reply [put]
func (h *ReviewHandler) ReplyToReview(c *gin.Context) {
	ctx := c.Request.Context()
	storeID := c.GetString("store_id")
	reviewID := c.Param("id")

	var req dto.SellerReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ResponseError(c, dto.NewError(http.StatusBadRequest, preference.ErrInvalidReqBody))
		return
	}

	err := h.reviewService.ReplyToReview(ctx, storeID, reviewID, req)
	if err != nil {
		helpers.ResponseError(c, err)
		return
	}
	helpers.ResponseSuccess(c, http.StatusOK, "Reply added successfully", nil)
}
