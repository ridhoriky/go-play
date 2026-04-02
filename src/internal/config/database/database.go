package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type DatabaseConfig struct {
	Enabled         bool          `yaml:"enabled"`
	Host            string        `yaml:"host"`
	Driver          string        `yaml:"driver"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Name            string        `yaml:"name"`
	SSLMode         string        `yaml:"ssl_mode"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
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
		log.Panic().Err(err).Msg(fmt.Sprintf("DB %s connection: FAILED", strings.ToUpper(cfgDb.Driver)))
		return nil, err
	}

	db.SetMaxOpenConns(cfgDb.MaxOpenConns)
	db.SetMaxIdleConns(cfgDb.MaxIdleConns)
	db.SetConnMaxLifetime(cfgDb.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfgDb.ConnMaxIdleTime)

	log.Debug().Msg(fmt.Sprintf("DB %s connection: SUCCESS", strings.ToUpper(cfgDb.Driver)))

	return db, nil
}
