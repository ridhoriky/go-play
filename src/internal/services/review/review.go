package review

import (
	"context"

	"ne-project/src/internal/models/dto"
)

type ReviewServiceItf interface {
	CreateReview(ctx context.Context, buyerID string, req dto.CreateReviewRequest) (*dto.ReviewResponse, error)
	GetProductReviews(ctx context.Context, productID string, params dto.GetReviewsQuery) (*dto.ReviewListResponse, error)
	ReplyToReview(ctx context.Context, storeID string, reviewID string, req dto.SellerReplyRequest) error
}
