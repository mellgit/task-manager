package cmd

import (
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/mellgit/task-manager/docs"
	"github.com/mellgit/task-manager/internal/auth"
	"github.com/mellgit/task-manager/internal/config"
	dbInit "github.com/mellgit/task-manager/internal/db"
	"github.com/mellgit/task-manager/internal/queue"
	"github.com/mellgit/task-manager/internal/task"
	"github.com/mellgit/task-manager/internal/worker"
	"github.com/mellgit/task-manager/pkg/logger"
	log "github.com/sirupsen/logrus"
)

// Up
// @title Task manager
// @version 1.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func Up() {

	cfgPath := flag.String("config", "config.yml", "config file path")
	flag.Parse()

	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.WithFields(log.Fields{
			"action": "config.LoadConfig",
		}).Fatal(err)
	}
	envCfg, err := config.LoadEnvConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"action": "config.LoadEnvConfig",
		}).Fatal(err)
	}

	if err = logger.SetUpLogger(*cfg); err != nil {
		log.WithFields(log.Fields{
			"action": "logger.SetUpLogger",
		}).Fatal(err)
	}

	log.Debugf("config: %+v", cfg)
	log.Debugf("env: %+v", envCfg)

	postgresClient, err := dbInit.PostgresClient(*envCfg)
	if err != nil {
		log.WithFields(log.Fields{
			"action": "dbInit.PostgresClient",
		}).Fatal(err)
	}

	kafkaAddr := fmt.Sprintf("%s:%d", envCfg.KafkaHost, envCfg.KafkaPort)
	producer, err := queue.NewProducer(kafkaAddr, envCfg.KafkaNameTopic, log.WithFields(log.Fields{"queue": "Producer"}))
	if err != nil {
		log.WithFields(log.Fields{
			"action": "queue.NewProducer",
		}).Fatal(err)
	}
	defer producer.Writer.Close()

	workerRepo := worker.NewRepo(postgresClient)
	workerService := worker.NewService(workerRepo, log.WithFields(log.Fields{"service": "Worker"}))
	consumer, err := queue.NewConsumer(kafkaAddr, envCfg.KafkaNameTopic, workerService, log.WithFields(log.Fields{"queue": "Consumer"}))
	if err != nil {
		log.WithFields(log.Fields{
			"action": "queue.NewConsumer",
		}).Fatal(err)
	}
	defer consumer.Reader.Close()

	go consumer.Start()

	app := fiber.New()
	{
		authRepo := auth.NewRepo(postgresClient)
		authService := auth.NewService(authRepo)
		authHandler := auth.NewHandler(authService, log.WithFields(log.Fields{"service": "AuthUser"}))
		authHandler.GroupHandler(app)

		taskRepo := task.NewRepo(postgresClient)
		taskService := task.NewService(taskRepo, producer)
		taskHandler := task.NewHandler(taskService, log.WithFields(log.Fields{"service": "Task"}))
		taskHandler.GroupHandler(app)

		app.Get("/swagger/*", swagger.HandlerDefault)

		log.Infof("http server listening %v:%v", envCfg.APIHost, envCfg.APIPort)
		log.WithFields(log.Fields{
			"action": "app.Listen",
		}).Fatal(app.Listen(fmt.Sprintf(":%v", envCfg.APIPort)))
	}

}
