package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"ne-project/src/internal/config/appconfig"
	"ne-project/src/internal/config/database"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

func main() {
	log.Println("Starting database seeder from CSV templates...")
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
	tables := []string{
		"carts", "wishlists", "reviews", "order_items", "orders",
		"transaction_details", "transactions", "product_images",
		"products", "categories", "stores", "users",
	}
	for _, t := range tables {
		if _, err := db.Exec("TRUNCATE TABLE " + t + " CASCADE"); err != nil {
			log.Printf("Error truncating %s: %v", t, err)
			return err
		}
	}
	log.Println("Truncated all tables")

	// Order of import matters due to foreign key constraints
	importOrder := []string{
		"users",
		"categories",
		"stores",
		"products",
		"product_images",
		"transactions",
		"transaction_details",
		"orders",
		"order_items",
		"reviews",
		"wishlists",
		"carts",
	}

	for _, table := range importOrder {
		filePath := fmt.Sprintf("db_templates/%s.csv", table)
		log.Printf("Importing data into %s from %s...", table, filePath)
		if err := seedFromCSV(db, table, filePath); err != nil {
			return err
		}
	}

	return nil
}

func seedFromCSV(db *sqlx.DB, tableName string, filePath string) error {
	//nolint:gosec // seeder script
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer func() { _ = file.Close() }()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading CSV %s: %w", filePath, err)
	}

	if len(records) < 2 {
		log.Printf("No data in %s, skipping", filePath)
		return nil
	}

	headers := records[0]
	columnsStr := strings.Join(headers, ", ")

	placeholders := make([]string, len(headers))
	for i := range headers {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columnsStr, strings.Join(placeholders, ", "))

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	for idx, record := range records[1:] {
		args := make([]any, len(record))
		for i, val := range record {
			args[i] = normalizeValue(tableName, headers[i], val)
		}

		_, err := tx.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("error inserting row %d into %s: %w", idx+2, tableName, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction for %s: %w", tableName, err)
	}

	log.Printf("Successfully imported %d rows into %s", len(records)-1, tableName)
	return nil
}

func normalizeValue(tableName string, header string, val string) any {
	if val == "" {
		return nil
	}

	if tableName == "transactions" && header == "status" && val == "failed" {
		return "canceled"
	}

	if tableName == "orders" && header == "status" {
		//nolint:misspell
		if val == "cancelled" {
			return "canceled"
		}
	}

	return val
}
