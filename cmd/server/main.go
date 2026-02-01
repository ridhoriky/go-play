package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/spf13/viper"

	"ne-project/internal/config"
	"ne-project/internal/database"
	"ne-project/internal/handlers"
	"ne-project/internal/repositories"
	"ne-project/internal/services"
)

func main() {

	// Load config
	if err := config.Load(); err != nil {
		log.Fatal("Failed load config:", err)
	}

	// Get DB config
	dbConn := viper.GetString("DB_CONN")

	if dbConn == "" {
		log.Fatal("DB_CONN is empty")
	}


	// Init DB
	db, err := database.InitDB(dbConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	defer db.Close()

	// Start server
	startServer(db)
}

func startServer(db *sql.DB) {

	// Dependency Injection
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Dependency Injection for Category
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Routing

	// Products
	http.HandleFunc("/api/products", productHandler.HandleProducts)
	http.HandleFunc("/api/products/", productHandler.HandleProductByID)

	// Categories
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)


	// Server config
	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port:", port)

	// Start HTTP Server
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
