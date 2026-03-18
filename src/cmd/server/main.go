package main

import (
	"fmt"
	"log"

	_ "ne-project/docs"
	"ne-project/src/internal/config"
	"ne-project/src/internal/database"
	"ne-project/src/internal/handlers/rest"
	"ne-project/src/internal/repositories"
	"ne-project/src/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

// @title           Kasir APP
// @version         1.0
// @description     KasirApp, a RESTful API built with Go for a simple cashier system using Gin and PostgreSQL.
// @host            localhost:8080
// @BasePath        /api/v1
func startServer(cfg *config.Config, db *sqlx.DB) {

	repo := repositories.NewRepository(db)
	service := services.NewServices(repo)
	handlers := rest.NewHandlers(service)

	// Routing
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Base Path Global
	v1 := r.Group("/api/v1")
	handlers.RegisterRoutes(v1)

	// Server config
	port := cfg.App.Port

	log.Println("Server running on port:", port)

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal("Failed to run server:", err)
	}
}
