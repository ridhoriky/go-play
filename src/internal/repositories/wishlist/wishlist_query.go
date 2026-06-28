package wishlist

const (
	queryAddWishlist = `
		INSERT INTO wishlists (buyer_id, product_id)
		VALUES (:buyer_id, :product_id)
		RETURNING id, buyer_id, product_id, created_at
	`

	queryRemoveWishlist = `
		DELETE FROM wishlists
		WHERE buyer_id = $1 AND product_id = $2
	`

	queryGetWishlistsByBuyerID = `
		SELECT w.id, w.buyer_id, w.product_id, w.created_at
		FROM wishlists w
		JOIN products p ON w.product_id = p.id
		WHERE w.buyer_id = $1 AND p.deleted_at IS NULL
		ORDER BY w.created_at DESC
		LIMIT $2 OFFSET $3
	`

	queryCountWishlistsByBuyerID = `
		SELECT COUNT(w.id)
		FROM wishlists w
		JOIN products p ON w.product_id = p.id
		WHERE w.buyer_id = $1 AND p.deleted_at IS NULL
	`

	queryCheckInWishlist = `
		SELECT EXISTS (
			SELECT 1 FROM wishlists
			WHERE buyer_id = $1 AND product_id = $2
		)
	`
)
