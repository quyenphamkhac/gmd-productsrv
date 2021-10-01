package driver

import (
	"log"
	"os"
	"sync"

	"github.com/streadway/amqp"
)

var (
	initRabbitmqOnce sync.Once
	rabbitmqConn     *amqp.Connection
)

func NewRabbitMQConn() *amqp.Connection {
	initRabbitmqOnce.Do(func() {
		var err error
		rabbitmqUrl := os.Getenv("RABBITMQ_URL")
		rabbitmqConn, err = amqp.Dial(rabbitmqUrl)
		if err != nil {
			log.Fatal(err)
		}
	})
	return rabbitmqConn
}
