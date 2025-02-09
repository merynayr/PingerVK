package ping

import (
	"context"
	"encoding/json"

	"github.com/merynayr/PingerVK/backend/internal/model"
	"github.com/merynayr/PingerVK/pkg/logger"
)

// Create - отправляет запрос в сервисный слой на создание данных о пингах и отправляет сообщение в Kafka
func (api *API) Create(message []byte) {
	var ping model.Pings

	// Привязываем JSON из запроса к модели Ping
	err := json.Unmarshal(message, &ping)
	if err != nil {
		logger.With("error", err).Error("Failed to unmarshal Kafka message")
		return
	}

	// Создаем пинг в сервисном слое
	_, err = api.pingService.Create(context.Background(), &ping)
	if err != nil {
		logger.Error("failed to create ping")
		return
	}
}
