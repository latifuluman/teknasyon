package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"listener/grpc"
	"listener/grpc/logs"
	"listener/grpc/mail"
)

// Consumer represents a RabbitMQ consumer.
type Consumer struct {
	conn *amqp.Connection // Connection to the RabbitMQ server.
}

// NewConsumer creates a new Consumer instance and sets up the RabbitMQ exchange.
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// setup initializes the RabbitMQ channel and declares the exchange.
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

// Payload represents the structure of the message payload.
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// MailPayload is the embedded type (in RequestPayload) that describes an email message to be sent.
type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// Listen starts listening for messages on the specified topics.
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		err = ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

// handlePayload processes the received payload based on its name.
func handlePayload(payload Payload) {
	switch payload.Name {
	case "log":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	case "mail":
		err := sendMail(payload)
		if err != nil {
			log.Println(err)
		}

	default:

	}
}

// sendMail sends an email by calling the mail-service.
func sendMail(msg Payload) error {
	var mreq mail.MailRequest
	json.Unmarshal([]byte(msg.Data), &mreq)
	c := grpc.GetGrpcMailClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := c.SendMail(ctx, &mreq)
	if err != nil {
		return err
	}
	return nil
}

// logEvent logs the event by calling the log-service.
func logEvent(entry Payload) error {
	c := grpc.GetGrpcLoggerClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: entry.Name,
			Data: entry.Data,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
