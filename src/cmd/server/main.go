package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

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
		log.Fatal().Err(err).Msg("Failed to connect to database")
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

	if err := serve(r, log, &cfg.App); err != nil {
		log.Fatal().Err(err).Msg("Server error")
	}
}

// serve starts the HTTP server and handles graceful shutdown on interrupt signals.
func serve(handler http.Handler, log *zerolog.Logger, appConfig *AppConfig) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", appConfig.Port),
		Handler:      handler,
		ReadTimeout:  appConfig.ReadTimeout,
		WriteTimeout: appConfig.WriteTimeout,
		IdleTimeout:  appConfig.IdleTimeout,
	}

	go func() {
		log.Info().Int("port", appConfig.Port).Msg("server starting " + appConfig.AppName + " version " + appConfig.Version + " in " + appConfig.Environment + " environment")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("unexpected server error")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("shutdown signal received, gracefully stopping...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), appConfig.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Info().Msg("server stopped gracefully")
	return nil
}
