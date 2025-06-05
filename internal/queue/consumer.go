package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(broker, topic string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{broker},
			Topic:     topic,
			Partition: 0,
			MaxBytes:  10e6, // 10MB
		}),
	}
}

func (c *Consumer) Start(process func(TaskPayload) error) {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("kafka read error: %v", err)
			continue
		}

		var payload TaskPayload
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			fmt.Printf("invalid task payload: %v", err)
			continue
		}

		fmt.Printf("processing task: %+v\n", payload)
		if err := process(payload); err != nil {
			fmt.Printf("task processing failed: %v", err)
		}
	}
}
