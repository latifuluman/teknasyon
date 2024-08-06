package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"os"
	"time"

	g "listener/grpc"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	RabbitMQ_URL       = "amqp://guest:guest@rabbitmq"
	MAIL_SERVICE_URL   = "mailer-service:50001"
	LOGGER_SERVICE_URL = "logger-service:50001"
)

func main() {

	// try to connect to rabbitmq
	rabbitConn, err := rabbitMQConnect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	mailConn := initGrpcMailConn()
	g.SetGrpcMailConn(mailConn)
	log.Printf("Connected to mail server")

	loggerConn := initGrpcLoggerConn()
	g.SetGrpcLoggerConn(loggerConn)
	log.Printf("Connected to logger server")

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func rabbitMQConnect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial(RabbitMQ_URL)
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}

func initGrpcLoggerConn() *grpc.ClientConn {

	for {
		conn, err := grpc.NewClient(LOGGER_SERVICE_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Printf("logger service connection error: %+v, waiting 5 sec...", err)
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		return conn
	}

}

func initGrpcMailConn() *grpc.ClientConn {

	for {
		conn, err := grpc.NewClient(MAIL_SERVICE_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("mail service connection error: %+v, waiting 5 sec...", err)
			time.Sleep(time.Duration(5) * time.Second)
			continue
		}
		return conn
	}

}
