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
	cfg           *config.Config
	logger        *log.Entry
	topics        []string
}

func NewConsumer(envCfg *config.EnvConfig, cfg *config.Config, serviceWorker worker.Service, logger *log.Entry) (*Consumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Producer.Retry.Max = 1
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.MaxMessageBytes = cfg.Broker.MaxMessageBytes * 1024 * 1024
	saramaConfig.ClientID = "sasl_scram_client"
	saramaConfig.Metadata.Full = true
	saramaConfig.Net.SASL.Enable = cfg.Broker.SASL.Enabled
	saramaConfig.Net.SASL.User = envCfg.KafkaUserName
	saramaConfig.Net.SASL.Password = envCfg.KafkaPassword
	saramaConfig.Net.SASL.Handshake = true

	if cfg.Broker.Algorithm == "sha512" {
		saramaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
		saramaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	} else if cfg.Broker.Algorithm == "sha256" {
		saramaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
		saramaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
	} else {
		logger.Fatalf("Invalid Broker Algorithm '%s'", cfg.Broker.Algorithm)
	}
	if cfg.Broker.KafkaTLS.Enabled {
		saramaConfig.Net.TLS.Enable = true
		saramaConfig.Net.TLS.Config = createTLSConfiguration(cfg)
	}

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
	// it is better, of course, to make a separate structure
	nc := &Consumer{serviceWorker: c.serviceWorker, logger: c.logger}
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

// Setup method before starting work
func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup a method for cleaning up resources after work is completed
func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim main method of message processing
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	type TaskPayload struct {
		TaskID uuid.UUID `json:"task_id"`
		UserID uuid.UUID `json:"user_id"`
	}
	for message := range claim.Messages() {

		c.logger.Infof("consume from kafka: topic=%v, value=%v", message.Topic, string(message.Value))
		//fmt.Printf("consume from kafka: topic=%v, value=%v", message.Topic, string(message.Value))
		data := new(TaskPayload)
		var _ = json.Unmarshal(message.Value, data)
		fromWorker := c.serviceWorker.GetPayload()
		fromWorker.UserID = data.UserID
		fromWorker.TaskID = data.TaskID
		if err := c.serviceWorker.Process(fromWorker); err != nil {
			c.logger.Errorf("failed to process task: %v\n", err)
			continue
		}
		// mark the message as processed
		session.MarkMessage(message, "")
	}

	return nil
}
