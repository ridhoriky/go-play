package review

const (
	insertReviewQuery = `
		INSERT INTO reviews (id, product_id, buyer_id, order_id, rating, comment, images, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	updateSellerReplyQuery = `
		UPDATE reviews 
		SET seller_reply = $1, seller_replied_at = NOW(), updated_at = NOW() 
		WHERE id = $2
	`

	getReviewByIDQuery = `
		SELECT id, product_id, buyer_id, order_id, rating, comment, images, seller_reply, seller_replied_at, created_at, updated_at
		FROM reviews 
		WHERE id = $1
	`

	getReviewsByProductIDBaseQuery = `
		SELECT r.id, r.product_id, r.buyer_id, r.order_id, r.rating, r.comment, r.images, r.seller_reply, r.seller_replied_at, r.created_at, r.updated_at,
		       u.name as buyer_name, u.avatar_url as buyer_avatar,
		       COUNT(*) OVER() as total_count
		FROM reviews r
		JOIN users u ON r.buyer_id = u.id
		WHERE r.product_id = $1
	`

	hasBuyerReviewedQuery = `
		SELECT EXISTS (
			SELECT 1 FROM reviews 
			WHERE buyer_id = $1 AND product_id = $2 AND order_id = $3
		)
	`

	getProductRatingSummaryQuery = `
		SELECT COALESCE(AVG(rating), 0) as avg, COUNT(id) as count
		FROM reviews
		WHERE product_id = $1
	`
)
