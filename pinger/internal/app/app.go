package app

import (
	"context"
	"flag"
	"sync"
	"time"

	"github.com/merynayr/PingerVK/pinger/internal/config"
	"github.com/merynayr/PingerVK/pinger/internal/service/ping"
	"github.com/merynayr/PingerVK/pkg/closer"
	"github.com/merynayr/PingerVK/pkg/logger"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

// App структура приложения
type App struct {
	serviceProvider *serviceProvider
}

// NewApp возвращает объект приложения
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run запускает приложение
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()
	address := a.serviceProvider.HTTPConfig().Address()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			ping.MonitorContainers(address)
			time.Sleep(10 * time.Second)
		}
	}()

	wg.Wait()

	return nil
}

// initDeps инициализирует все зависимости
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initLogger,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initConfig(_ context.Context) error {
	flag.Parse()
	err := config.Load(configPath)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.Init(a.serviceProvider.LoggerConfig().Level())

	return nil
}
