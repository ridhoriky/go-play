package main

import (
	"database/sql"
	"fmt"
	"log"
	"ne-project/internal/config"
	"ne-project/internal/database"
	"ne-project/internal/handlers"
	"ne-project/internal/repositories"
	"ne-project/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}


	// Init DB
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	defer db.Close()

	// Start server
	startServer(cfg, db)
}

func startServer(cfg *config.Config, db *sql.DB) {
	
	repo := repositories.NewRepository(db)
	service := services.NewServices(repo)
	handlers := handlers.NewHandlers(service)

	// Routing
	r := gin.Default()
	handlers.RegisterRoutes(r)

	// Server config
	port := cfg.App.Port

	log.Println("Server running on port:", port)

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
