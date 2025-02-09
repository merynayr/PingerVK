package consumer

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/merynayr/PingerVK/backend/internal/client/kafka"
)

var (
	errConsumerNotInitialized = errors.New("consumer is not initialized")
)

// kafkaConsumerGroup структура для потребления сообщений через Consumer Group
type kafkaConsumerGroup struct {
	group   sarama.ConsumerGroup
	handler Handler
}

// Handler реализует интерфейс sarama.ConsumerGroupHandler
type Handler struct {
	handlerFunc func(*sarama.ConsumerMessage) error
}

// New возвращает новый консьюмер Kafka с Consumer Group
func New(brokersAddresses []string, groupID string, cfg *sarama.Config) (kafka.Consumer, error) {
	group, err := sarama.NewConsumerGroup(brokersAddresses, groupID, cfg)
	if err != nil {
		return nil, err
	}
	return &kafkaConsumerGroup{group: group}, nil
}

// Consume запускает консьюмера в фоновом режиме
func (c *kafkaConsumerGroup) Consume(ctx context.Context, topics []string, handler func(*sarama.ConsumerMessage) error) error {
	if c.group == nil {
		return errConsumerNotInitialized
	}

	c.handler = Handler{handlerFunc: handler}

	for {
		err := c.group.Consume(ctx, topics, &c.handler)
		if err != nil {
			return fmt.Errorf("failed to consume messages: %w", err)
		}
	}
}

// Close закрывает соединение консьюмера
func (c *kafkaConsumerGroup) Close() error {
	if c.group == nil {
		return errConsumerNotInitialized
	}
	return c.group.Close()
}

// Setup вызывается перед стартом потребления партиций
func (h *Handler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup вызывается после завершения работы с партициями
func (h *Handler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim обрабатывает сообщения из партиции
func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		err := h.handlerFunc(msg)
		if err != nil {
			log.Printf("Ошибка обработки сообщения: %v", err)
			continue
		}

		// Явно коммитим offset после успешной обработки
		session.MarkMessage(msg, "")
	}
	return nil
}
