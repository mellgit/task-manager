package queue

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/mellgit/task-manager/internal/worker"

	log "github.com/sirupsen/logrus"
)

type Consumer struct {
	Reader        sarama.Consumer
	topic         string
	serviceWorker worker.Service
	logger        *log.Entry
}

func NewConsumer(brokerAddress, topic string, serviceWorker worker.Service, logger *log.Entry) (*Consumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{brokerAddress}, saramaConfig)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		Reader:        consumer,
		topic:         topic,
		serviceWorker: serviceWorker,
		logger:        logger,
	}, nil
}

func (c *Consumer) Start() {

	consumer, err := c.Reader.ConsumePartition(c.topic, 0, sarama.OffsetOldest)
	if err != nil {
		c.logger.Errorf("Error creating consumer for topic %s: %s", c.topic, err)
	}
	defer consumer.Close()

	log.Info("starting consumer server")

	for {
		select {
		case err := <-consumer.Errors():
			c.logger.Errorf("fail to consume from kafka %v\n", err)
		case msg := <-consumer.Messages():
			c.logger.Infof("consume from kafka %s", msg.Value)
			data := new(TaskPayload)
			if err := json.Unmarshal(msg.Value, data); err != nil {
				log.Errorf("failed to unmarshal message: %v\n", err)
				continue
			}
			// taskPayload in the queue package and in worker are the same
			// but golang considers them different
			// so it is necessary to convert the types
			fromWorker := c.serviceWorker.GetPayload()
			fromWorker.UserID = data.UserID
			fromWorker.TaskID = data.TaskID
			if err := c.serviceWorker.Process(fromWorker); err != nil {
            	c.logger.Errorf("failed to process task: %v\n", err)
				continue
			}

		}
	}
}
