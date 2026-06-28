package product_image

const (
	getProductImagesQuery = `
		SELECT id, product_id, url, alt_text, sort_order, is_primary, created_at
		FROM product_images
		WHERE product_id = $1
		ORDER BY sort_order ASC, created_at ASC
	`
	insertProductImageQuery = `
		INSERT INTO product_images (product_id, url, alt_text, sort_order, is_primary)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	deleteProductImageQuery = `
		DELETE FROM product_images WHERE id = $1
	`
	setPrimaryImageQuery = `
		UPDATE product_images SET is_primary = true WHERE id = $1 AND product_id = $2
	`
	unsetPrimaryImageQuery = `
		UPDATE product_images SET is_primary = false WHERE product_id = $1
	`
	countProductImagesQuery = `
		SELECT COUNT(*) FROM product_images WHERE product_id = $1
	`
)
