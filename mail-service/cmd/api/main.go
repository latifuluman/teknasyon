package main

import (
	"log"
)

// Config holds the application's configuration.
type Config struct {
	Mailer Mail
}

func main() {
	// Initialize the application configuration with mail settings.
	app := Config{
		Mailer: createMail(),
	}

	// Start the gRPC server and listen for connections.
	err := app.gRPCListen()
	if err != nil {
		log.Panicf("gRPC listen failed with err: %+v\n", err)
	}
}
