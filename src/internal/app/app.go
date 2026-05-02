package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"ne-project/src/internal/config/appconfig"
	"ne-project/src/internal/config/database"
	"ne-project/src/internal/config/logger"
	"ne-project/src/internal/config/middleware"
	"ne-project/src/internal/config/token"
	"ne-project/src/internal/handlers/rest/routes"
	"ne-project/src/internal/repositories"
	"ne-project/src/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(cfg *appconfig.Config) {
	log := logger.InitLogger(cfg.Logger)

	db, err := database.InitDB(*log, &cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("app - Run - failed to initialize database")
	}
	defer closeDatabase(log, db)

	tokenSvc, err := token.InitToken(*log, cfg.Token)
	if err != nil {
		log.Err(err).Msg("app - Run - failed to initialize token service")
		return
	}

	router := setupRouter(log, cfg, db, tokenSvc)

	srv := initServer(cfg, router)
	waitForShutdown(log, srv, cfg.App)
}

func closeDatabase(log *zerolog.Logger, db *sqlx.DB) {
	if db == nil {
		return
	}
	log.Info().Str("component", "database").Msg("closing database connection")
	if err := db.Close(); err != nil {
		log.Error().Err(err).Str("component", "database").Msg("failed to close database connection")
	}
}

func setupRouter(log *zerolog.Logger, cfg *appconfig.Config, db *sqlx.DB, tokenSvc *token.Token) *gin.Engine {
	mw := middleware.InitMiddleware(*log, tokenSvc, cfg.RateLimit.Limit, cfg.RateLimit.Window)
	repo := repositories.NewRepository(db)
	service := services.NewServices(repo, tokenSvc, db)
	handlers := routes.NewHandlers(db, service)

	r := gin.New()
	r.Use(mw.Logger(), mw.CORS(), gin.Recovery())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	handlers.RegisterRoutes(v1, tokenSvc, *mw)

	return r
}

func initServer(cfg *appconfig.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      handler,
		ReadTimeout:  cfg.App.ReadTimeout,
		WriteTimeout: cfg.App.WriteTimeout,
		IdleTimeout:  cfg.App.IdleTimeout,
	}
}

func waitForShutdown(log *zerolog.Logger, srv *http.Server, appCfg appconfig.AppConfig) {
	serverError := make(chan error, 1)

	go func() {
		log.Info().
			Str("app", appCfg.AppName).
			Str("version", appCfg.Version).
			Str("env", appCfg.Environment).
			Int("port", appCfg.Port).
			Msg("HTTP server starting")

		serverError <- srv.ListenAndServe()
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info().Str("signal", s.String()).Msg("app - Run - signal received")
	case err := <-serverError:
		if err != http.ErrServerClosed {
			log.Error().Err(err).Msg("app - Run - unexpected server error")
		}
	}

	log.Info().Msg("shutting down HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), appCfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("app - Run - server forced to shutdown")
	}

	log.Info().Msg("server stopped gracefully")
}
