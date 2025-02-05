package repository

import (
	"context"

	"github.com/merynayr/PingerVK/backend/internal/model"
)

// PingRepository - интерфейс репо слоя ping
type PingRepository interface {
	Create(ctx context.Context, ping *model.Pings) (int64, error)
	Get(ctx context.Context) ([]*model.Pings, error)
}
