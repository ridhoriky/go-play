package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"ne-project/src/internal/config/appconfig"
	"ne-project/src/internal/config/database"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

// makeUUID generates a deterministic valid UUID based on a hex prefix and an index.
func makeUUID(group string, idx int) string {
	return fmt.Sprintf("01918341-%4s-0000-0000-%012d", group, idx)
}

func main() {
	log.Println("Starting database seeder...")

	cfg, err := appconfig.LoadConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	db, err := database.InitDB(&logger, &cfg.Database)
	if err != nil {
		log.Fatalf("DB connection error: %s", err)
	}
	defer func() { _ = db.Close() }()

	if err := seed(db); err != nil {
		log.Printf("Seeding error: %s", err)
	} else {
		log.Println("Seeding completed successfully")
	}
}

func seed(db *sqlx.DB) error {
	total := 50
	passwordHash := "$2a$12$Z0HjJ.Zp6T0311NqM1VzUeqQ4B0M1/o/2D1T0.eSg4U78s.7K1Nue" //nolint:gosec

	log.Printf("Seeding %d records per table...", total)

	if err := seedUsers(db, total, passwordHash); err != nil {
		return err
	}
	if err := seedStores(db, total); err != nil {
		return err
	}
	if err := seedCategories(db, total); err != nil {
		return err
	}
	if err := seedProducts(db, total); err != nil {
		return err
	}
	if err := seedProductImages(db, total); err != nil {
		return err
	}
	if err := seedTransactions(db, total); err != nil {
		return err
	}
	if err := seedTransactionDetails(db, total); err != nil {
		return err
	}
	if err := seedOrders(db, total); err != nil {
		return err
	}
	if err := seedOrderItems(db, total); err != nil {
		return err
	}
	if err := seedReviews(db, total); err != nil {
		return err
	}
	if err := seedWishlists(db, total); err != nil {
		return err
	}
	if err := seedCarts(db, total); err != nil {
		return err
	}

	return nil
}

func seedUsers(db *sqlx.DB, total int, passwordHash string) error {
	var users []string
	for i := 1; i <= total; i++ {
		id := makeUUID("a000", i)
		role := "user"
		switch i % 3 {
		case 1:
			role = "admin"
		case 2:
			role = "seller"
		}
		users = append(users, fmt.Sprintf("('%s', 'user%d@example.com', '%s', 'User %d', '%s', true, true)", id, i, passwordHash, i, role))
	}
	return execBatch(db, "users", "id, email, password_hash, name, role, is_active, is_verified", users, "id")
}

func seedStores(db *sqlx.DB, total int) error {
	var stores []string
	for i := 1; i <= total; i++ {
		id := makeUUID("b000", i)
		userId := makeUUID("a000", i)
		stores = append(stores, fmt.Sprintf("('%s', '%s', 'Store %d', 'store-%d', 'Description for store %d', true)", id, userId, i, i, i))
	}
	return execBatch(db, "stores", "id, user_id, store_name, slug, description, is_verified", stores, "id")
}

func seedCategories(db *sqlx.DB, total int) error {
	var categories []string
	for i := 1; i <= total; i++ {
		id := makeUUID("c000", i)
		categories = append(categories, fmt.Sprintf("('%s', 'Category %d', 'Description %d')", id, i, i))
	}
	return execBatch(db, "categories", "id, name, description", categories, "id")
}

func seedProducts(db *sqlx.DB, total int) error {
	var products []string
	for i := 1; i <= total; i++ {
		id := makeUUID("d000", i)
		catId := makeUUID("c000", (i%total)+1)
		storeId := makeUUID("b000", (i%total)+1)
		price := (rand.Intn(100) + 1) * 10000 //nolint:gosec
		stock := rand.Intn(100) + 10          //nolint:gosec
		products = append(products, fmt.Sprintf("('%s', 'Product %d', %d, %d, '%s', '%s', 'Desc %d', 'product-%d')", id, i, price, stock, catId, storeId, i, i))
	}
	return execBatch(db, "products", "id, name, price, stock, category_id, store_id, description, slug", products, "id")
}

func seedProductImages(db *sqlx.DB, total int) error {
	var productImages []string
	for i := 1; i <= total; i++ {
		id := makeUUID("e000", i)
		prodId := makeUUID("d000", i)
		productImages = append(productImages, fmt.Sprintf("('%s', '%s', 'https://via.placeholder.com/500', 'Alt %d', true)", id, prodId, i))
	}
	return execBatch(db, "product_images", "id, product_id, url, alt_text, is_primary", productImages, "id")
}

func seedTransactions(db *sqlx.DB, total int) error {
	var transactions []string
	for i := 1; i <= total; i++ {
		id := makeUUID("f000", i)
		transactions = append(transactions, fmt.Sprintf("('%s', 150000, 'completed')", id))
	}
	return execBatch(db, "transactions", "id, total_amount, status", transactions, "id")
}

func seedTransactionDetails(db *sqlx.DB, total int) error {
	var transactionDetails []string
	for i := 1; i <= total; i++ {
		id := makeUUID("a100", i)
		transId := makeUUID("f000", i)
		prodId := makeUUID("d000", i)
		transactionDetails = append(transactionDetails, fmt.Sprintf("('%s', '%s', '%s', 'Product %d', 1, 150000, 150000)", id, prodId, transId, i))
	}
	return execBatch(db, "transaction_details", "id, product_id, transaction_id, product_name, quantity, price, subtotal", transactionDetails, "id")
}

func seedOrders(db *sqlx.DB, total int) error {
	var orders []string
	for i := 1; i <= total; i++ {
		id := makeUUID("b100", i)
		buyerId := makeUUID("a000", (i%total)+1)
		storeId := makeUUID("b000", (i%total)+1)
		orders = append(orders, fmt.Sprintf("('%s', '%s', '%s', 'ORD-%04d', 50000, 'delivered', '{\"address\":\"Street %d\"}', 10000, 'transfer')", id, buyerId, storeId, i, i))
	}
	return execBatch(db, "orders", "id, buyer_id, store_id, order_number, total_amount, status, shipping_address, shipping_cost, payment_method", orders, "id")
}

func seedOrderItems(db *sqlx.DB, total int) error {
	var orderItems []string
	for i := 1; i <= total; i++ {
		id := makeUUID("c100", i)
		orderId := makeUUID("b100", i)
		prodId := makeUUID("d000", i)
		orderItems = append(orderItems, fmt.Sprintf("('%s', '%s', '%s', 'Product %d', 1, 50000, 50000)", id, orderId, prodId, i))
	}
	return execBatch(db, "order_items", "id, order_id, product_id, product_name, quantity, price, subtotal", orderItems, "id")
}

func seedReviews(db *sqlx.DB, total int) error {
	var reviews []string
	for i := 1; i <= total; i++ {
		id := makeUUID("d100", i)
		prodId := makeUUID("d000", i)
		buyerId := makeUUID("a000", (i%total)+1)
		orderId := makeUUID("b100", i)
		reviews = append(reviews, fmt.Sprintf("('%s', '%s', '%s', '%s', 5, 'Great %d')", id, prodId, buyerId, orderId, i))
	}
	return execBatch(db, "reviews", "id, product_id, buyer_id, order_id, rating, comment", reviews, "id")
}

func seedWishlists(db *sqlx.DB, total int) error {
	var wishlists []string
	for i := 1; i <= total; i++ {
		id := makeUUID("e100", i)
		buyerId := makeUUID("a000", (i%total)+1)
		prodId := makeUUID("d000", i)
		wishlists = append(wishlists, fmt.Sprintf("('%s', '%s', '%s')", id, buyerId, prodId))
	}
	return execBatch(db, "wishlists", "id, buyer_id, product_id", wishlists, "id")
}

func seedCarts(db *sqlx.DB, total int) error {
	var carts []string
	for i := 1; i <= total; i++ {
		id := makeUUID("f100", i)
		buyerId := makeUUID("a000", (i%total)+1)
		prodId := makeUUID("d000", i)
		carts = append(carts, fmt.Sprintf("('%s', '%s', '%s', 1)", id, buyerId, prodId))
	}
	return execBatch(db, "carts", "id, buyer_id, product_id, quantity", carts, "id")
}

func execBatch(db *sqlx.DB, table string, columns string, values []string, conflictCol string) error {
	if len(values) == 0 {
		return nil
	}

	// Process in chunks to avoid any query length limits (though PG can handle quite a lot)
	chunkSize := 100
	for i := 0; i < len(values); i += chunkSize {
		end := min(i+chunkSize, len(values))

		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s ON CONFLICT (%s) DO NOTHING",
			table, columns, strings.Join(values[i:end], ", "), conflictCol)

		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error inserting into %s: %v", table, err)
			return err
		}
	}

	log.Printf("Seeded %s successfully", table)
	return nil
}
