package env

import (
	"os"

	"github.com/pkg/errors"

	"github.com/merynayr/PingerVK/pinger/internal/config"
)

const (
	backendAPIEnvName = "BACKEND_API_URL"
)

type httpConfig struct {
	address string
}

// NewHTTPConfig returns new http-server config
func NewHTTPConfig() (config.HTTPConfig, error) {
	address := os.Getenv(backendAPIEnvName)
	if len(address) == 0 {
		return nil, errors.New("http port not found")
	}

	return &httpConfig{
		address: address,
	}, nil
}

func (cfg *httpConfig) Address() string {
	return cfg.address
}
