package ping

import (
	"context"

	"github.com/merynayr/PingerVK/backend/internal/model"
)

func (s *srv) Create(ctx context.Context, ping *model.Pings) (int64, error) {
	ID, err := s.pingRepository.Create(ctx, ping)
	if err != nil {
		return 0, err
	}

	return ID, nil
}
