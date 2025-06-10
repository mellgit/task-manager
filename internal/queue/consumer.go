package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/mellgit/task-manager/internal/config"
	"github.com/mellgit/task-manager/internal/worker"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Consumer struct {
	Reader        sarama.ConsumerGroup
	serviceWorker worker.Service
	envCfg        *config.EnvConfig
	logger        *log.Entry
	topics        []string
}

func NewConsumer(envCfg *config.EnvConfig, serviceWorker worker.Service, logger *log.Entry) (*Consumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true
	brokers := strings.Split(envCfg.KafkaBrokers, ",")
	topics := strings.Split(envCfg.KafkaTopic, ",")
	groupID := envCfg.KafkaGroup
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("consumer group initialization error: %v", err)
	}

	return &Consumer{
		Reader:        consumerGroup,
		topics:        topics,
		serviceWorker: serviceWorker,
		logger:        logger,
	}, nil
}

func (c *Consumer) Start() {

	ctx := context.Background()
	nc := &MyConsumer{serviceWorker: c.serviceWorker}
	for {
		err := c.Reader.Consume(ctx, c.topics, nc)
		if err != nil {
			c.logger.Errorf("error when processing a consumer group: %v", err)
		}
		select {
		case <-ctx.Done():
			panic(ctx.Err())
		default:
			c.logger.Info("default")
		}
	}

}

// MyConsumer Структура обработчика группы потребителей
type MyConsumer struct {
	serviceWorker worker.Service
}

// Setup Метод настройки перед началом работы
func (m *MyConsumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup Метод очистки ресурсов после завершения работы
func (m *MyConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim Основной метод обработки сообщений
func (m *MyConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	type TaskPayload struct {
		TaskID uuid.UUID `json:"task_id"`
		UserID uuid.UUID `json:"user_id"`
	}
	for message := range claim.Messages() {

		//m.worker.Logger.Infof("consume from kafka: topic=%v, value=%v", message.Topic, string(message.Value))

		data := new(TaskPayload)
		var _ = json.Unmarshal(message.Value, data)
		fromWorker := m.serviceWorker.GetPayload()
		fromWorker.UserID = data.UserID
		fromWorker.TaskID = data.TaskID
		if err := m.serviceWorker.Process(fromWorker); err != nil {
			//c.logger.Errorf("failed to process task: %v\n", err)
			continue
		}
		// Отмечаем сообщение как обработанное
		session.MarkMessage(message, "")
	}

	return nil
}
