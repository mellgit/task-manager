package queue

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"time"

	"github.com/segmentio/kafka-go"
)

type TaskPayload struct {
	TaskID uuid.UUID `json:"task_id"`
	UserID uuid.UUID `json:"user_id"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokerAddress, topicName string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokerAddress),
			Topic:    topicName,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Publish(payload TaskPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(time.Now().String()),
		Value: data,
	}

	return p.writer.WriteMessages(context.Background(), msg)
}
