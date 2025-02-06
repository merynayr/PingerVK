package ping

import (
	"context"

	"github.com/merynayr/PingerVK/backend/internal/model"
)

func (s *srv) Create(ctx context.Context, ping *model.Pings) (int64, error) {
	var ID int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		var exist bool

		exist, errTx = s.pingRepository.ContainerExists(ctx, ping.ID)
		if errTx != nil {
			return errTx
		}
		if !exist {
			ID, errTx = s.pingRepository.Create(ctx, ping)
		} else {
			errTx = s.pingRepository.UpdateInfo(ctx, ping)
		}
		if errTx != nil {
			return errTx
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return ID, nil
}
