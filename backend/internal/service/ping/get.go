package ping

import (
	"context"

	"github.com/merynayr/PingerVK/backend/internal/model"
)

func (s *srv) Get(ctx context.Context) ([]*model.GetPings, error) {
	pings, err := s.pingRepository.Get(ctx)
	if err != nil {
		return nil, err
	}

	return pings, nil
}
