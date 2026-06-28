package dto

import (
	"time"
)

type CreateReviewRequest struct {
	ProductID string   `json:"product_id" binding:"required,uuid"`
	OrderID   string   `json:"order_id" binding:"required,uuid"`
	Rating    int      `json:"rating" binding:"required,min=1,max=5"`
	Comment   *string  `json:"comment"`
	Images    []string `json:"images"`
}

type SellerReplyRequest struct {
	Reply string `json:"reply" binding:"required"`
}

type ReviewResponse struct {
	ID              string     `json:"id"`
	BuyerName       string     `json:"buyer_name,omitempty"`
	BuyerAvatar     *string    `json:"buyer_avatar,omitempty"`
	Rating          int        `json:"rating"`
	Comment         *string    `json:"comment,omitempty"`
	Images          []string   `json:"images"`
	SellerReply     *string    `json:"seller_reply,omitempty"`
	SellerRepliedAt *time.Time `json:"seller_replied_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type GetReviewsQuery struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

type ReviewListResponse struct {
	Items         []ReviewResponse `json:"items"`
	AverageRating float64          `json:"average_rating"`
	TotalReviews  int              `json:"total_reviews"`
	Meta          PaginationMeta   `json:"meta"`
}
