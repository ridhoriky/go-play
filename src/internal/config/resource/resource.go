package resource

import (
	"ne-project/src/internal/config/appconfig"
	"ne-project/src/internal/config/database"
	redisconfig "ne-project/src/internal/config/redis"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Resources struct {
	DB    *sqlx.DB
	Redis *redis.Client
}

func InitResources(log *zerolog.Logger, cfg *appconfig.Config) (*Resources, error) {
	db, err := database.InitDB(log, &cfg.Database)
	if err != nil {
		return nil, err
	}

	redisClient, err := redisconfig.InitRedis(log, &cfg.Redis)
	if err != nil {
		log.Info().Str("component", "database").Msg("closing database connection due to redis failure")
		if err := db.Close(); err != nil {
			log.Error().Err(err).Str("component", "database").Msg("failed to close database connection")
		}
		return nil, err
	}

	return &Resources{
		DB:    db,
		Redis: redisClient,
	}, nil
}

func (r *Resources) Close(log *zerolog.Logger) {
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
