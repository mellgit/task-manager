package queue

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/mellgit/task-manager/internal/config"
	log "github.com/sirupsen/logrus"
	"strings"
)

type TaskPayload struct {
	TaskID uuid.UUID `json:"task_id"`
	UserID uuid.UUID `json:"user_id"`
}

type Producer struct {
	Writer sarama.SyncProducer
	envCfg *config.EnvConfig
	logger *log.Entry
	topic  string
}

func NewProducer(envCfg *config.EnvConfig, logger *log.Entry) (*Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5

	brokers := strings.Split(envCfg.KafkaBrokers, ",")
	producer, err := sarama.NewSyncProducer(brokers, saramaConfig)
	if err != nil {
		return nil, err
	}
	return &Producer{
		Writer: producer,
		topic:  envCfg.KafkaTopic,
		logger: logger,
	}, nil
}

func (p *Producer) Publish(payload TaskPayload) error {

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshalling task payload: %s", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := p.Writer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("error sending message: %s", err)
	}
	p.logger.Infof("message sent to partition %d at offset %d\n", partition, offset)
	return nil
}
