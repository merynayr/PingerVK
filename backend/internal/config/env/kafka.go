package env

import (
	"errors"
	"os"
	"strings"

	"github.com/IBM/sarama"

	"github.com/merynayr/PingerVK/backend/internal/config"
)

const (
	kafkaBrokersEnvName     = "KAFKA_BROKERS"
	kafkaTopicsEnvName      = "KAFKA_TOPICS"
	kafkaGroupEnvName       = "KAFKA_GROUP"
	consumerRetryMax        = 5
	consumerReturnSuccesses = true
)

type kafkaConsumerConfig struct {
	brokers []string
	topics  []string
	group   string
}

// NewKafkaConsumerConfig возвращает конфиг kafka
func NewKafkaConsumerConfig() (config.KafkaConsumerConfig, error) {
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

	group := os.Getenv(kafkaGroupEnvName)
	if len(group) == 0 {
		return nil, errors.New("kafka group not found")
	}

	return &kafkaConsumerConfig{
		brokers: brokers,
		topics:  topics,
		group:   group,
	}, nil
}

// Brokers возвращает список адрессов брокеров
func (cfg *kafkaConsumerConfig) Brokers() []string {
	if cfg.brokers == nil {
		return []string{}
	}

	return cfg.brokers
}

// Topics возвращает список топиков
func (cfg *kafkaConsumerConfig) Topics() []string {
	if cfg.topics == nil {
		return []string{}
	}

	return cfg.topics
}

// Group возвращает группу Kafka
func (cfg *kafkaConsumerConfig) Group() string {
	return cfg.group
}

// Config returns sarama consumer config
func (cfg *kafkaConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Idempotent = true
	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Net.MaxOpenRequests = 1
	config.Producer.Retry.Max = consumerRetryMax
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Return.Successes = consumerReturnSuccesses
	return config
}
