package main

import (
	"fmt"

	_ "ne-project/docs"
	"ne-project/src/internal/config/database"
	"ne-project/src/internal/config/logger"
	"ne-project/src/internal/handlers/rest"
	"ne-project/src/internal/handlers/rest/system"
	"ne-project/src/internal/repositories"
	"ne-project/src/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Kasir APP
// @version         1.0
// @description     KasirApp, a RESTful API built with Go for a simple cashier system using Gin and PostgreSQL.
// @host            localhost:8080
// @BasePath        /api/v1
func main() {

	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	log := logger.InitLogger(cfg.Logger)

	// Init DB
	db, err := database.InitDB(log, &cfg.Database)
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("Failed to initialize database: %s", err))
	}

	defer db.Close()

	repo := repositories.NewRepository(db)
	service := services.NewServices(repo)
	handlers := rest.NewHandlers(service)

	// Routing
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Base Path Global
	v1 := r.Group("/api/v1")
	handlers.RegisterRoutes(v1)

	// System Routes
	system.NewSystemHandler(db).RegisterRoutes(r)

	// Server config
	port := cfg.App.Port

	log.Info().Msg(fmt.Sprintf("Server running on port: %d", port))

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal().Msg(fmt.Sprintf("Failed to run server: %s", err))
	}
}
