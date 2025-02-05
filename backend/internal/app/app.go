package app

import (
	"context"
	"flag"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/merynayr/PingerVK/backend/internal/config"
	"github.com/merynayr/PingerVK/pkg/closer"
	"github.com/merynayr/PingerVK/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

// App структура приложения
type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
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

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
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
		a.initHTTPServer,
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

func (a *App) initHTTPServer(ctx context.Context) error {
	logger.Init(a.serviceProvider.LoggerConfig().Level())

	router := gin.Default()

	a.serviceProvider.PingAPI(ctx)

	a.serviceProvider.pingAPI.RegisterRoutes(router)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.HTTPConfig().Address(),
		Handler:           corsMiddleware.Handler(router),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return nil
}

func (a *App) runHTTPServer() error {
	log := logger.With("address", a.serviceProvider.HTTPConfig().Address())
	log.Info("HTTP server is running")

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
