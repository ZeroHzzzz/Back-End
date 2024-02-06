package models

import "github.com/streadway/amqp"

type RabbitMQMiddleware struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}
