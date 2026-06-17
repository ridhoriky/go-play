package resource

import (
	"context"

	"ne-project/src/internal/config/appconfig"
	"ne-project/src/internal/config/database"
	"ne-project/src/internal/config/metrics"
	redisconfig "ne-project/src/internal/config/redis"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/sdk/metric"
)

type Resources struct {
	DB    *sqlx.DB
	Redis *redis.Client
	MP    *metric.MeterProvider
}

func InitResources(log *zerolog.Logger, cfg *appconfig.Config) (*Resources, error) {
	db, err := database.InitDB(log, &cfg.Database)
	if err != nil {
		return nil, err
	}

	redisClient, err := redisconfig.InitRedis(log, &cfg.Redis)
	if err != nil {
		log.Info().Str("component", "database").Msg("closing database connection due to redis failure")
		if closeErr := db.Close(); closeErr != nil {
			log.Error().Err(err).Str("component", "database").Msg("failed to close database connection")
		}
		return nil, err
	}

	mp := metrics.InitMetrics()

	return &Resources{
		DB:    db,
		Redis: redisClient,
		MP:    mp,
	}, nil
}

func (r *Resources) Close(log *zerolog.Logger) {
	if r.MP != nil {
		log.Info().Str("component", "metrics").Msg("closing Metrics Provider connection")
		if err := r.MP.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Str("component", "metrics").Msg("failed to close Metrics Provider connection")
		}
	}

	if r.Redis != nil {
		log.Info().Str("component", "redis").Msg("closing Redis connection")
		if err := r.Redis.Close(); err != nil {
			log.Error().Err(err).Str("component", "redis").Msg("failed to close Redis connection")
		}
	}

	if r.DB != nil {
		log.Info().Str("component", "database").Msg("closing database connection")
		if err := r.DB.Close(); err != nil {
			log.Error().Err(err).Str("component", "database").Msg("failed to close database connection")
		}
	}
}
