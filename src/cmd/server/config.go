package main

import (
	"fmt"
	"ne-project/src/internal/config/database"
	"ne-project/src/internal/config/logger"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig               `yaml:"app"`
	Database database.DatabaseConfig `yaml:"database"`
	Logger   logger.LoggerOptions    `yaml:"logger"`
}

type AppConfig struct {
	Port int `yaml:"port"`
}

func LoadConfig() (*Config, error) {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	viper.SetEnvKeyReplacer(
		strings.NewReplacer(".", "_"),
	)

	viper.AutomaticEnv()

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	cfg.App.Port = parseInt(viper.GetString("APP.PORT"), 8080)
	cfg.Database.Port = parseInt(viper.GetString("DATABASE.PORT"), 6543)

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

	if val := os.Getenv("PORT"); val != "" {
		cfg.App.Port = parseInt(val, cfg.App.Port)
	}

}

func parseInt(s string, defaultVal int) int {
	var val int
	if _, err := fmt.Sscanf(s, "%d", &val); err == nil {
		return val
	}

	return defaultVal
}
