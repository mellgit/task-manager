package queue

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type TaskPayload struct {
	TaskID uuid.UUID `json:"task_id"`
	UserID uuid.UUID `json:"user_id"`
}

type Producer struct {
	Writer sarama.SyncProducer
	topic  string
	logger *log.Entry
}

func NewProducer(brokerAddress, topic string, logger *log.Entry) (*Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer([]string{brokerAddress}, saramaConfig)
	if err != nil {
		return nil, err
	}
	return &Producer{
		Writer: producer,
		topic:  topic,
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
