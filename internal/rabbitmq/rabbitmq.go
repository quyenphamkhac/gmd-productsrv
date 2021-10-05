package rabbitmq

import (
	"fmt"
	"log"

	"github.com/quyenphamkhac/gmd-productsrv/config"
	"github.com/streadway/amqp"
)

func NewRabbitMQConn(cfg *config.Config) (*amqp.Connection, error) {
	connStr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)
	log.Println("rabbitmq connection string: ", connStr)
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Printf("unable connect rabbit mq, %v", err)
		return nil, err
	}
	return conn, nil
}
