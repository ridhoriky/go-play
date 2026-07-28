package redis

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type RedisOptions struct {
	Addr     string        `yaml:"addr" env:"REDIS_ADDR" env-default:"localhost:6379"`
	Username string        `yaml:"username" env:"REDIS_USERNAME" env-default:""`
	Password string        `yaml:"password" env:"REDIS_PASSWORD" env-default:""`
	DB       int           `yaml:"db" env:"REDIS_DB" env-default:"0"`
	PoolSize int           `yaml:"pool_size" env:"REDIS_POOL_SIZE" env-default:"10"`
	Timeout  time.Duration `yaml:"timeout" env:"REDIS_TIMEOUT" env-default:"5s"`
	TLS      bool          `yaml:"tls" env:"REDIS_TLS" env-default:"false"`
}

func InitRedis(log *zerolog.Logger, cfg *RedisOptions) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:      cfg.Addr,
		Username:  cfg.Username,
		Password:  cfg.Password,
		DB:        cfg.DB,
		PoolSize:  cfg.PoolSize,
		TLSConfig: &tls.Config{},
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Info().Str("component", "redis").Msg("successfully connected to Redis")
	return client, nil
}
