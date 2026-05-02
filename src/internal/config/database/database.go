package database

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type DatabaseConfig struct {
	Enabled         bool          `yaml:"enabled" env:"DB_ENABLED" env-default:"false"`
	Host            string        `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Driver          string        `yaml:"driver" env:"DB_DRIVER" env-default:"postgres"`
	Port            int           `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User            string        `yaml:"user" env:"DB_USER" env-default:"postgres"`
	Password        string        `yaml:"password" env:"DB_PASSWORD" env-default:"mypassword"`
	Name            string        `yaml:"name" env:"DB_NAME" env-default:"kasir_db"`
	SSLMode         string        `yaml:"ssl_mode" env:"DB_SSL_MODE" env-default:"disable"`
	MaxOpenConns    int           `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS" env-default:"25"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS" env-default:"5"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME" env-default:"1h"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" env:"DB_CONN_MAX_IDLE_TIME" env-default:"30m"`
}

func GetConnectionString(db *DatabaseConfig) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.SSLMode,
	)
}

func InitDB(log zerolog.Logger, cfgDb *DatabaseConfig) (*sqlx.DB, error) {

	db, err := sqlx.Connect(cfgDb.Driver, GetConnectionString(cfgDb))
	if err != nil {
		log.Error().Err(err).Msgf("Failed to connect to database: %s", err.Error())
		return nil, err
	}

	db.SetMaxOpenConns(cfgDb.MaxOpenConns)
	db.SetMaxIdleConns(cfgDb.MaxIdleConns)
	db.SetConnMaxLifetime(cfgDb.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfgDb.ConnMaxIdleTime)

	log.Info().Msg("Database connection established successfully")

	return db, nil
}
