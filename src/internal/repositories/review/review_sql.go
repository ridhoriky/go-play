package review

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"ne-project/src/internal/models/dto"
	"ne-project/src/internal/models/entity"

	"github.com/jmoiron/sqlx"
)

type reviewRepository struct {
	db *sqlx.DB
}

func NewReviewRepository(db *sqlx.DB) ReviewRepositoryItf {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *entity.Review) error {
	_, err := r.db.ExecContext(ctx, insertReviewQuery,
		review.ID, review.ProductID, review.BuyerID, review.OrderID, review.Rating, review.Comment, review.Images, review.CreatedAt, review.UpdatedAt,
	)
	return err
}

func (r *reviewRepository) GetByProductID(ctx context.Context, productID string, params *dto.GetReviewsQuery) ([]entity.ReviewWithBuyer, int, error) {
	query := getReviewsByProductIDBaseQuery
	args := []any{productID}

	// Add sorting by newest
	query += " ORDER BY r.created_at DESC"

	if params.Limit > 0 {
		query += " LIMIT $2"
		args = append(args, params.Limit)

		if params.Page > 0 {
			offset := (params.Page - 1) * params.Limit
			query += " OFFSET $3"
			args = append(args, offset)
		}
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var reviews []entity.ReviewWithBuyer
	var total int

	for rows.Next() {
		var rev entity.ReviewWithBuyer
		var rowCount int

		if scanErr := rows.Scan(
			&rev.ID, &rev.ProductID, &rev.BuyerID, &rev.OrderID, &rev.Rating, &rev.Comment, &rev.Images, &rev.SellerReply, &rev.SellerRepliedAt, &rev.CreatedAt, &rev.UpdatedAt,
			&rev.BuyerName, &rev.BuyerAvatar,
			&rowCount,
		); scanErr != nil {
			return nil, 0, scanErr
		}

		reviews = append(reviews, rev)
		if total == 0 {
			total = rowCount
		}
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

func (r *reviewRepository) GetByID(ctx context.Context, id string) (*entity.Review, error) {
	var rev entity.Review
	err := r.db.QueryRowContext(ctx, getReviewByIDQuery, id).Scan(
		&rev.ID, &rev.ProductID, &rev.BuyerID, &rev.OrderID, &rev.Rating, &rev.Comment, &rev.Images, &rev.SellerReply, &rev.SellerRepliedAt, &rev.CreatedAt, &rev.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, dto.NewError(http.StatusNotFound, "review not found")
	}
	if err != nil {
		return nil, err
	}
	return &rev, nil
}

func (r *reviewRepository) HasBuyerReviewed(ctx context.Context, buyerID string, productID string, orderID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, hasBuyerReviewedQuery, buyerID, productID, orderID).Scan(&exists)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	return exists, nil
}

func (r *reviewRepository) AddSellerReply(ctx context.Context, id string, reply string) error {
	res, err := r.db.ExecContext(ctx, updateSellerReplyQuery, reply, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return dto.NewError(http.StatusNotFound, "review not found")
	}
	return nil
}

func (r *reviewRepository) GetProductRatingSummary(ctx context.Context, productID string) (float64, int, error) {
	var avg float64
	var count int
	err := r.db.QueryRowContext(ctx, getProductRatingSummaryQuery, productID).Scan(&avg, &count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, nil
		}
		return 0, 0, err
	}
	return avg, count, nil
}
