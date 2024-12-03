package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"trinity/configs"
	"trinity/db/psql"
	"trinity/db/redis_client"
	"trinity/internal/model"
	"trinity/pkg/app"
	httpserver "trinity/pkg/http_server"
)

var (
	configMap map[string]string
	cfg       configs.Config
	appConfig app.Config
)

func initializeLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel)
}

func loadConfigs() {
	var err error
	configMap, err = configs.GetConfigFromEnv()
	if err != nil {
		log.Fatalf("Failed to load configs: %v", err)
	}

	cfg = configs.NewConfig(configMap)
	appConfig = app.NewAppConfig(cfg)
}

func initializeDBConnections() {
	manager := psql.NewDBManager()
	err := manager.Connect("db1", configMap["dbAddr"])
	if err != nil {
		log.Fatalf("Failed to connect to db1: %v", err)
	}
	model.InitUserDB(manager)
	model.InitCampaignDB(manager)
	model.InitVoucherDB(manager)
}

func initializeRedisConnection() {
	err := redis_client.InitializeRedis(configMap["redisAddr"], "", 0)
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}
}

func init() {
	initializeLogger()
	loadConfigs()
	initializeDBConnections()
	initializeRedisConnection()
}

// @title API Documentation
// @version 1.0
// @description Develop a promotional campaign system for the Trinity app, enabling a 30% discount on Silver subscription plans for the first 100 users registering via campaign links. The system will generate time-limited vouchers to ensure efficient campaign management and user engagement.
// @host localhost:9888
// @BasePath /
func main() {
	server := httpserver.NewServer(&appConfig, &cfg)
	server.StartWithGracefulShutdown()
}
