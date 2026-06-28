package review

import (
	"context"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
)

type ReviewRepositoryItf interface {
	Create(ctx context.Context, review *entity.Review) error
	GetByProductID(ctx context.Context, productID string, params *dto.GetReviewsQuery) ([]entity.ReviewWithBuyer, int, error)
	GetByID(ctx context.Context, id string) (*entity.Review, error)
	HasBuyerReviewed(ctx context.Context, buyerID string, productID string, orderID string) (bool, error)
	AddSellerReply(ctx context.Context, id string, reply string) error
	GetProductRatingSummary(ctx context.Context, productID string) (float64, int, error)
}
