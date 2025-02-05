package app

import (
	"context"
	"log"

	"github.com/merynayr/PingerVK/backend/internal/config"
	"github.com/merynayr/PingerVK/backend/internal/config/env"
	"github.com/merynayr/PingerVK/backend/internal/repository"
	"github.com/merynayr/PingerVK/backend/internal/service"
	"github.com/merynayr/PingerVK/pkg/client/db"
	"github.com/merynayr/PingerVK/pkg/client/db/pg"
	"github.com/merynayr/PingerVK/pkg/client/db/transaction"
	"github.com/merynayr/PingerVK/pkg/closer"

	pingAPI "github.com/merynayr/PingerVK/backend/internal/api/ping"

	pingService "github.com/merynayr/PingerVK/backend/internal/service/ping"

	pingRepository "github.com/merynayr/PingerVK/backend/internal/repository/ping"
)

// Структура приложения со всеми зависимости
type serviceProvider struct {
	pgConfig     config.PGConfig
	httpConfig   config.HTTPConfig
	loggerConfig config.LoggerConfig

	dbClient  db.Client
	txManager db.TxManager

	pingRepository repository.PingRepository
	pingService    service.PingService
	pingAPI        *pingAPI.API
}

// NewServiceProvider возвращает новый объект API слоя
func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}
		s.pgConfig = cfg
	}
	return s.pgConfig
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

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) PingRepository(ctx context.Context) repository.PingRepository {
	if s.pingRepository == nil {
		s.pingRepository = pingRepository.NewRepository(s.DBClient(ctx))
	}

	return s.pingRepository
}

func (s *serviceProvider) PingService(ctx context.Context) service.PingService {
	if s.pingService == nil {
		s.pingService = pingService.NewService(
			s.PingRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.pingService
}

func (s *serviceProvider) PingAPI(ctx context.Context) *pingAPI.API {
	if s.pingAPI == nil {
		s.pingAPI = pingAPI.NewAPI(s.PingService(ctx))
	}

	return s.pingAPI
}
