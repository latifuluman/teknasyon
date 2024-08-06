package main

import (
	"context"
	"log"
	"log-service/data"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"                    // Port for the web server.
	mongoURL = "mongodb://mongo:27017" // URL for connecting to MongoDB.
	gRpcPort = "50001"                 // Port for the gRPC server.
)

var client *mongo.Client // Global MongoDB client.

// Config holds the application's configuration including database models.
type Config struct {
	Models data.Models // Models for interacting with the database.
}

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// Create a context with timeout for disconnecting from MongoDB.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Close MongoDB connection when the application exits.
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// Initialize the application configuration with database models.
	app := Config{
		Models: data.New(client),
	}

	// Start the gRPC server and listen for connections.
	err = app.gRPCListen()
	if err != nil {
		log.Panicf("gRPC listen failed with err: %+v\n", err)
	}
}

// connectToMongo establishes a connection to MongoDB and returns the client.
func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to mongo!")

	return c, nil
}
