package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// declareExchange declares an exchange on the given RabbitMQ channel.
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name of the exchange
		"topic",      // type of the exchange (topic)
		true,         // durable? (survives broker restarts)
		false,        // auto-deleted? (deleted when no longer in use)
		false,        // internal? (used internally by RabbitMQ)
		false,        // no-wait? (do not wait for the server response)
		nil,          // arguments? (additional options for the exchange)
	)
}

// declareRandomQueue declares a random, exclusive queue on the given RabbitMQ channel.
func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name? (empty string creates a random name)
		false, // durable? (queue does not survive broker restarts)
		false, // delete when unused? (queue deleted when no consumers)
		true,  // exclusive? (queue only accessible by the connection that declares it)
		false, // no-wait? (do not wait for the server response)
		nil,   // arguments? (additional options for the queue)
	)
}
