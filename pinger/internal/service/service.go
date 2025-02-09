package service

// PingService интерфейс сервисного слоя user
type PingService interface {
	SendContainer(pingTopicName string) error
}
