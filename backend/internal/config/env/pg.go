package env

import (
	"errors"
	"os"

	"github.com/merynayr/PingerVK/backend/internal/config"
)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

// NewPGConfig returns new postgres config
func NewPGConfig() (config.PGConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("pg dsn not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}

// DSN returns a full database's connection string
func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
