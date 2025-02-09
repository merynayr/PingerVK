package env

import (
	"errors"
	"os"
	"strings"

	"github.com/IBM/sarama"
	"github.com/merynayr/PingerVK/pinger/internal/config"
)

const (
	kafkaBrokersEnvName     = "KAFKA_BROKERS"
	kafkaTopicsEnvName      = "KAFKA_TOPICS"
	producerRetryMax        = 5
	producerReturnSuccesses = true
)

type kafkaProducerConfig struct {
	brokers []string
	topics  []string
}

// NewKafkaProducerConfig возвращает конфиг kafka
func NewKafkaProducerConfig() (config.KafkaProducerConfig, error) {
	brokersStr := os.Getenv(kafkaBrokersEnvName)
	if len(brokersStr) == 0 {
		return nil, errors.New("kafka brokers address not found")
	}
	brokers := strings.Split(brokersStr, ",")

	topicsStr := os.Getenv(kafkaTopicsEnvName)
	if len(topicsStr) == 0 {
		return nil, errors.New("kafka topics not found")
	}
	topics := strings.Split(topicsStr, ",")

	return &kafkaProducerConfig{
		brokers: brokers,
		topics:  topics,
	}, nil
}

// Brokers возвращает список адрессов брокеров
func (cfg *kafkaProducerConfig) Brokers() []string {
	if cfg.brokers == nil {
		return []string{}
	}

	return cfg.brokers
}

// Topics возвращает список топиков
func (cfg *kafkaProducerConfig) Topics() []string {
	if cfg.topics == nil {
		return []string{}
	}

	return cfg.topics
}

// Config returns sarama producer config
func (cfg *kafkaProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = producerRetryMax
	config.Producer.Return.Successes = producerReturnSuccesses

	return config
}
