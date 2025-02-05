package config

import (
	"github.com/joho/godotenv"
)

// Load читает .env файл по указанному пути
// и загружает переменные в проект
func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

// PGConfig is interface of a postgres config
type PGConfig interface {
	DSN() string
}

// HTTPConfig is interface of a http config
type HTTPConfig interface {
	Address() string
}

// LoggerConfig интерфейс конфига логгера
type LoggerConfig interface {
	Level() string
}
