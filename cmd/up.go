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
	"github.com/mellgit/task-manager/internal/pkg/logger"
	log "github.com/sirupsen/logrus"
)

// Up
// @title Task manager
// @version 1.0
// @host localhost:3000
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

	// TODO kafka init

	app := fiber.New()
	{
		authRepo := auth.NewRepo(postgresClient)
		authService := auth.NewService(authRepo)
		authHandler := auth.NewHandler(authService, log.WithFields(log.Fields{"service": "AuthUser"}))
		authHandler.GroupHandler(app)

		app.Get("/swagger/*", swagger.HandlerDefault)

		log.Infof("http server listening %v:%v", envCfg.APIHost, envCfg.APIPort)
		log.WithFields(log.Fields{
			"action": "app.Listen",
		}).Fatal(app.Listen(fmt.Sprintf(":%v", envCfg.APIPort)))
	}

}
