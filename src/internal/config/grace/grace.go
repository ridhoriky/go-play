package grace

import (
	"context"
	"ne-project/src/internal/config/appconfig"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
)

func WaitForShutdown(log *zerolog.Logger, srv *http.Server, appCfg appconfig.AppConfig) {
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
