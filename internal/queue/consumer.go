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

func NewConsumer(broker, groupID, topic string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			GroupID: groupID,
			Topic:   topic,
		}),
	}
}

func (c *Consumer) Start(process func(TaskPayload) error) {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Kafka read error: %v", err)
			continue
		}

		var payload TaskPayload
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			fmt.Printf("Invalid task payload: %v", err)
			continue
		}

		fmt.Printf("Processing task: %+v\n", payload)
		if err := process(payload); err != nil {
			fmt.Printf("Task processing failed: %v", err)
		}
	}
}
