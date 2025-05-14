package cmd

import (
	"flag"
	//"github.com/gofiber/fiber/v2"

	"github.com/mellgit/task-manager/internal/config"

	//dbInit "github.com/mellgit/task-manager/internal/db"
	"github.com/mellgit/task-manager/internal/pkg/logger"
	log "github.com/sirupsen/logrus"
)

func Up() {

	//db init
	// kafka init
	//fiber init

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

	//postgresClient, err := dbInit.PostgresClient(*envCfg)
	//if err != nil {
	//	log.WithFields(log.Fields{
	//		"action": "dbInit.PostgresClient",
	//	}).Fatal(err)
	//}
	//
	//app := fiber.New()
	//
	//{
	//
	//}

}
