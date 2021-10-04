package adapter

import "github.com/streadway/amqp"

type productEventConsumer struct {
	amqpConn *amqp.Connection
}

func NewProductEventConsumer(amqpConn *amqp.Connection) *productEventConsumer {
	return &productEventConsumer{
		amqpConn: amqpConn,
	}
}
