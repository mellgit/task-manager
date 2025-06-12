package queue

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/mellgit/task-manager/internal/config"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type TaskPayload struct {
	TaskID uuid.UUID `json:"task_id"`
	UserID uuid.UUID `json:"user_id"`
}

type Producer struct {
	Writer sarama.SyncProducer
	envCfg *config.EnvConfig
	cfg    *config.Config
	logger *log.Entry
	topic  string
}

func NewProducer(envCfg *config.EnvConfig, cfg *config.Config, logger *log.Entry) (*Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 5
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

func createTLSConfiguration(cfg *config.Config) (t *tls.Config) {
	t = &tls.Config{
		InsecureSkipVerify: cfg.TLSSkipVerify,
	}

	if cfg.CertFile != "" && cfg.KeyFile != "" && cfg.CAFile != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			log.Fatal(err)
		}

		caCert, err := os.ReadFile(cfg.CAFile)
		if err != nil {
			log.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: cfg.TLSSkipVerify,
		}
	}
	return t
}
