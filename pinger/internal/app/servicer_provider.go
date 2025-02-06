package app

import (
	"log"

	"github.com/merynayr/PingerVK/pinger/internal/config"
	"github.com/merynayr/PingerVK/pinger/internal/config/env"
)

// Структура приложения со всеми зависимости
type serviceProvider struct {
	httpConfig   config.HTTPConfig
	loggerConfig config.LoggerConfig
}

// NewServiceProvider возвращает новый объект API слоя
func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) LoggerConfig() config.LoggerConfig {
	if s.loggerConfig == nil {
		cfg, err := env.NewLoggerConfig()
		if err != nil {
			log.Fatalf("failed to get logger config:%v", err)
		}

		s.loggerConfig = cfg
	}

	return s.loggerConfig
}
