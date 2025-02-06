package service

import (
	"context"

	"github.com/merynayr/PingerVK/backend/internal/model"
)

// PingService интерфейс сервисного слоя user
type PingService interface {
	Create(ctx context.Context, ping *model.Pings) (int64, error)
	Get(ctx context.Context) ([]*model.GetPings, error)
}
