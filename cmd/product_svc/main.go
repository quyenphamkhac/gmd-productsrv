package main

import (
	"log"

	"github.com/quyenphamkhac/gmd-productsrv/config"
	"github.com/quyenphamkhac/gmd-productsrv/internal/logger"
	"github.com/quyenphamkhac/gmd-productsrv/internal/rabbitmq"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/server"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	svcLogger, err := logger.NewServiceLogger(cfg)
	svcLogger.InitLogger()
	if err != nil {
		log.Fatalf("unable init logger: %v", err)
	}

	rabbitmqConn, err := rabbitmq.NewRabbitMQConn(cfg)
	if err != nil {
		svcLogger.Fatalf("connect rabbitmq: %v", err)
	}
	defer rabbitmqConn.Close()

	app := server.NewProductServer(svcLogger, cfg, rabbitmqConn)
	app.Run()
}
