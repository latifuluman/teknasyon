package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Emitter represents an event emitter for RabbitMQ.
type Emitter struct {
	connection *amqp.Connection // Connection to the RabbitMQ server.
}

// setup initializes the RabbitMQ channel and declares the exchange.
func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

// Push sends an event to the RabbitMQ exchange with a specified severity.
func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel")

	// Publish the event to the RabbitMQ exchange.
	err = channel.Publish(
		"logs_topic", // Exchange name
		severity,     // Routing key (severity)
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event), // Event data
		},
	)
	if err != nil {
		return err
	}
	log.Printf("Pushed to RabbitMQ")
	return nil
}

// NewEventEmitter creates a new Emitter instance and sets up the RabbitMQ exchange.
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
