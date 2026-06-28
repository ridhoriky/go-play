package review

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"time"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"
	"ne-project/src/internal/repositories/order"
	"ne-project/src/internal/repositories/product"
	"ne-project/src/internal/repositories/review"

	"github.com/google/uuid"
)

type reviewService struct {
	reviewRepo  review.ReviewRepositoryItf
	orderRepo   order.OrderRepositoryItf
	productRepo product.ProductRepositoryItf
}

func NewReviewService(
	reviewRepo review.ReviewRepositoryItf,
	orderRepo order.OrderRepositoryItf,
	productRepo product.ProductRepositoryItf,
) ReviewServiceItf {
	return &reviewService{
		reviewRepo:  reviewRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *reviewService) CreateReview(ctx context.Context, buyerID string, req dto.CreateReviewRequest) (*dto.ReviewResponse, error) {
	// Validate: order exists & belongs to buyer
	ord, items, err := s.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, err
	}
	if ord.BuyerID != buyerID {
		return nil, dto.NewError(http.StatusForbidden, "you can only review your own orders")
	}

	// Validate: order status is "delivered"
	if ord.Status != "delivered" {
		return nil, dto.NewError(http.StatusBadRequest, "order must be delivered before reviewing")
	}

	// Validate product is in order
	productInOrder := false
	for i := range items {
		if items[i].ProductID == req.ProductID {
			productInOrder = true
			break
		}
	}
	if !productInOrder {
		return nil, dto.NewError(http.StatusBadRequest, "product not found in this order")
	}

	// Validate: buyer hasn't reviewed this product for this order
	hasReviewed, err := s.reviewRepo.HasBuyerReviewed(ctx, buyerID, req.ProductID, req.OrderID)
	if err != nil {
		return nil, err
	}
	if hasReviewed {
		return nil, dto.NewError(http.StatusConflict, "you have already reviewed this product for this order")
	}

	// Create review
	imagesJSON, err := json.Marshal(req.Images)
	if err != nil {
		return nil, dto.NewError(http.StatusInternalServerError, "failed to marshal images")
	}
	if len(req.Images) == 0 {
		imagesJSON = []byte("[]")
	}

	now := time.Now()
	rev := &entity.Review{
		ID:        uuid.NewString(),
		ProductID: req.ProductID,
		BuyerID:   buyerID,
		OrderID:   req.OrderID,
		Rating:    req.Rating,
		Comment:   req.Comment,
		Images:    imagesJSON,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.reviewRepo.Create(ctx, rev)
	if err != nil {
		return nil, err
	}

	// Note: Product and Store rating_avg triggers will automatically update via DB

	res := &dto.ReviewResponse{
		ID:        rev.ID,
		Rating:    rev.Rating,
		Comment:   rev.Comment,
		Images:    req.Images,
		CreatedAt: rev.CreatedAt,
	}

	return res, nil
}

func (s *reviewService) GetProductReviews(ctx context.Context, productID string, params dto.GetReviewsQuery) (*dto.ReviewListResponse, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	reviews, total, err := s.reviewRepo.GetByProductID(ctx, productID, &params)
	if err != nil {
		return nil, err
	}

	avgRating, totalReviews, err := s.reviewRepo.GetProductRatingSummary(ctx, productID)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ReviewResponse, 0)
	for i := range reviews {
		r := &reviews[i]
		var imgs []string
		_ = json.Unmarshal(r.Images, &imgs)

		items = append(items, dto.ReviewResponse{
			ID:              r.ID,
			BuyerName:       r.BuyerName,
			BuyerAvatar:     r.BuyerAvatar,
			Rating:          r.Rating,
			Comment:         r.Comment,
			Images:          imgs,
			SellerReply:     r.SellerReply,
			SellerRepliedAt: r.SellerRepliedAt,
			CreatedAt:       r.CreatedAt,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return &dto.ReviewListResponse{
		Items:         items,
		AverageRating: avgRating,
		TotalReviews:  totalReviews,
		Meta: dto.PaginationMeta{
			Total:      total,
			Page:       params.Page,
			Limit:      params.Limit,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *reviewService) ReplyToReview(ctx context.Context, storeID string, reviewID string, req dto.SellerReplyRequest) error {
	rev, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return err
	}

	prodDetail, err := s.productRepo.GetByID(ctx, rev.ProductID)
	if err != nil {
		return err
	}
	if prodDetail == nil {
		return dto.NewError(http.StatusNotFound, "product not found")
	}

	if prodDetail.StoreID != storeID {
		return dto.NewError(http.StatusForbidden, "product does not belong to your store")
	}

	return s.reviewRepo.AddSellerReply(ctx, reviewID, req.Reply)
}
