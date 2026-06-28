package cart

const (
	addCartQuery = `
		INSERT INTO carts (id, buyer_id, product_id, quantity, created_at, updated_at)
		VALUES (:id, :buyer_id, :product_id, :quantity, :created_at, :updated_at)
	`

	getCartByBuyerIDQuery = `
		SELECT id, buyer_id, product_id, quantity, created_at, updated_at
		FROM carts
		WHERE buyer_id = $1
		ORDER BY created_at DESC
	`

	getCartByIDQuery = `
		SELECT id, buyer_id, product_id, quantity, created_at, updated_at
		FROM carts
		WHERE id = $1
	`

	updateCartQuantityQuery = `
		UPDATE carts
		SET quantity = $1, updated_at = NOW()
		WHERE id = $2
	`

	deleteCartQuery = `
		DELETE FROM carts
		WHERE id = $1
	`

	deleteCartByBuyerIDQuery = `
		DELETE FROM carts
		WHERE buyer_id = $1
	`

	getCartByBuyerAndProductQuery = `
		SELECT id, buyer_id, product_id, quantity, created_at, updated_at
		FROM carts
		WHERE buyer_id = $1 AND product_id = $2
	`
)
