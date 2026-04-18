package main

import (
	"fmt"
	"ne-project/src/internal/config/database"
	"ne-project/src/internal/config/logger"
	"ne-project/src/internal/config/token"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	App      AppConfig               `yaml:"app"`
	Database database.DatabaseConfig `yaml:"database"`
	Logger   logger.LoggerOptions    `yaml:"logger"`
	Token    token.TokenOptions      `yaml:"logger"`
}

type AppConfig struct {
	Port int `yaml:"port"`
}

func LoadConfig() (*Config, error) {

	configPath := "config.yaml"

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	overrideWithEnv(&cfg)

	return &cfg, nil
}

func overrideWithEnv(cfg *Config) {
	if val := os.Getenv("DB_HOST"); val != "" {
		cfg.Database.Host = val
	}

	if val := os.Getenv("DB_PORT"); val != "" {
		cfg.Database.Port = parseInt(val, cfg.Database.Port)
	}

	if val := os.Getenv("DB_USER"); val != "" {
		cfg.Database.User = val
	}

	if val := os.Getenv("DB_PASSWORD"); val != "" {
		cfg.Database.Password = val
	}

	if val := os.Getenv("DB_NAME"); val != "" {
		cfg.Database.Name = val
	}

	if val := os.Getenv("SECRET_ACCESS_TOKEN"); val != "" {
		cfg.Token.SecretAccessToken = []byte(val)
	}

	if val := os.Getenv("SECRET_REFRESH_TOKEN"); val != "" {
		cfg.Token.SecretRefreshToken = []byte(val)
	}

}

func parseInt(s string, defaultVal int) int {
	var val int
	if _, err := fmt.Sscanf(s, "%d", &val); err == nil {
		return val
	}

	return defaultVal
}
