package ping

import (
	"github.com/merynayr/PingerVK/backend/internal/repository"
	"github.com/merynayr/PingerVK/backend/internal/service"
	"github.com/merynayr/PingerVK/pkg/client/db"
)

// Структура сервисного слоя с объектами репо слоя
// и транзакционного менеджера
type srv struct {
	pingRepository repository.PingRepository
	txManager      db.TxManager
}

// NewService возвращает объект сервисного слоя
func NewService(
	pingRepository repository.PingRepository,
	txManager db.TxManager,
) service.PingService {
	return &srv{
		pingRepository: pingRepository,
		txManager:      txManager,
	}
}
