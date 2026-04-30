package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "ne-project/docs"
	"ne-project/src/internal/config/database"
	"ne-project/src/internal/config/logger"
	"ne-project/src/internal/config/middleware"
	"ne-project/src/internal/config/token"
	"ne-project/src/internal/handlers/rest"
	"ne-project/src/internal/handlers/rest/system"
	"ne-project/src/internal/repositories"
	"ne-project/src/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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
	db, err := database.InitDB(*log, &cfg.Database)
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("Failed to initialize database: %s", err))
	}

	defer db.Close()

	// Token Service Initialization
	tokenSvc := token.InitToken(*log, cfg.Token)

	// Middleware Initialization
	mw := middleware.InitMiddleware(*log, tokenSvc, cfg.RateLimit.Limit, cfg.RateLimit.Window)

	repo := repositories.NewRepository(db)
	service := services.NewServices(repo, tokenSvc, db)
	handlers := rest.NewHandlers(service)
	r := gin.New()
	r.Use(mw.Logger(), mw.CORS(), gin.Recovery())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Base Path Global
	v1 := r.Group("/api/v1")
	handlers.RegisterRoutes(v1, tokenSvc, *mw)

	// System Routes
	system.NewSystemHandler(db).RegisterRoutes(v1)

	// Server config
	port := cfg.App.Port

	log.Info().Msg(fmt.Sprintf("Server running on port: %d", port))

	if err := serve(r, port, log); err != nil {
		log.Fatal().Err(err).Msg("Server error")
	}
}

// serve menjalankan HTTP server dengan Gin dan menangani graceful shutdown.
func serve(handler http.Handler, port int, log *zerolog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Int("port", port).Msg("server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("unexpected server error")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("shutdown signal received, gracefully stopping...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Info().Msg("server stopped gracefully")
	return nil
}
