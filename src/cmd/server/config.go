package main

import (
	"log"
	"ne-project/src/internal/config/database"
	"ne-project/src/internal/config/logger"
	"ne-project/src/internal/config/middleware"
	"ne-project/src/internal/config/token"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App       AppConfig                     `yaml:"app"`
	Database  database.DatabaseConfig       `yaml:"database"`
	Logger    logger.LoggerOptions          `yaml:"logger"`
	Token     token.TokenOptions            `yaml:"token"`
	RateLimit middleware.RateLimiterOptions `yaml:"rate_limit"`
}

type AppConfig struct {
	AppName         string        `yaml:"name" env:"APP_NAME" env-default:"myapp"`
	Version         string        `yaml:"version" env:"APP_VERSION" env-default:"1.0.0"`
	Port            int           `yaml:"port" env:"APP_PORT" env-default:"8080"`
	Environment     string        `yaml:"environment" env:"APP_ENVIRONMENT" env-default:"development"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env:"APP_WRITE_TIMEOUT" env-default:"15s"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env:"APP_READ_TIMEOUT" env-default:"15s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env:"APP_IDLE_TIMEOUT" env-default:"60s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"APP_SHUTDOWN_TIMEOUT" env-default:"30s"`
}

func LoadConfig() (*Config, error) {

	var cfg Config
	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Config loaded: %+v\n", cfg)

	return &cfg, nil
}
