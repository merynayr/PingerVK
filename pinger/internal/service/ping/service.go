package ping

import (
	"github.com/merynayr/PingerVK/pinger/internal/client/kafka"
	"github.com/merynayr/PingerVK/pinger/internal/service"
)

// Структура сервисного слоя
type srv struct {
	kafkaProducer kafka.Producer
}

// NewService возвращает объект сервисного слоя
func NewService(
	kafkaProducer kafka.Producer,
) service.PingService {
	return &srv{
		kafkaProducer: kafkaProducer,
	}
}
