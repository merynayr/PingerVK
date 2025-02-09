package app

import (
	"context"
	"log"

	"github.com/merynayr/PingerVK/pinger/internal/client/kafka"
	"github.com/merynayr/PingerVK/pinger/internal/client/kafka/producer"
	"github.com/merynayr/PingerVK/pinger/internal/config"
	"github.com/merynayr/PingerVK/pinger/internal/config/env"
	"github.com/merynayr/PingerVK/pinger/internal/service"
	pingSrv "github.com/merynayr/PingerVK/pinger/internal/service/ping"
	"github.com/merynayr/PingerVK/pkg/closer"
)

// Структура приложения со всеми зависимости
type serviceProvider struct {
	httpConfig          config.HTTPConfig
	loggerConfig        config.LoggerConfig
	kafkaProducerConfig config.KafkaProducerConfig

	pingService service.PingService

	kafkaProducer kafka.Producer
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

func (s *serviceProvider) KafkaProducerConfig() config.KafkaProducerConfig {
	if s.kafkaProducerConfig == nil {
		cfg, err := env.NewKafkaProducerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka producer config: %s", err)
		}

		s.kafkaProducerConfig = cfg
	}

	return s.kafkaProducerConfig
}

func (s *serviceProvider) KafkaProducer(_ context.Context) kafka.Producer {
	if s.kafkaProducer == nil {
		p, err := producer.New(s.KafkaProducerConfig().Brokers(), s.KafkaProducerConfig().Config())
		if err != nil {
			log.Fatalf("failed to create kafka producer: %v", err)
		}

		closer.Add(p.Close)
		s.kafkaProducer = p
	}

	return s.kafkaProducer
}

func (s *serviceProvider) PingService(ctx context.Context) service.PingService {
	if nil == s.pingService {
		s.pingService = pingSrv.NewService(s.KafkaProducer(ctx))
	}
	return s.pingService
}
