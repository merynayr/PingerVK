package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

// Consumer интерфейс продюсера кафки
type Consumer interface {
	Consume(ctx context.Context, topics []string, handler func(*sarama.ConsumerMessage) error) error
	Close() error
}
