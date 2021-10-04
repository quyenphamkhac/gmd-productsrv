package adapter

import "github.com/streadway/amqp"

type productEventPublisher struct {
	amqpConn *amqp.Connection
}

func NewProductEventPublisher(amqpConn *amqp.Connection) *productEventPublisher {
	return &productEventPublisher{
		amqpConn: amqpConn,
	}
}
