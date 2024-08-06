package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"

	"google.golang.org/grpc"
)

// LogServer is a gRPC server for handling log-related requests.
type LogServer struct {
	logs.UnimplementedLogServiceServer             // Embeds the unimplemented server for forward compatibility.
	Models                             data.Models // Models for interacting with the database.
}

// WriteLog handles the incoming log writing requests via gRPC.
func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// Create a log entry from the request data.
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	// Insert the log entry into the database.
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	// Return success response.
	res := &logs.LogResponse{Result: "logged!"}
	return res, nil
}

// gRPCListen starts the gRPC server and listens for incoming connections.
func (app *Config) gRPCListen() error {
	// Create a TCP listener on the specified port.
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		return err
	}

	// Create a new gRPC server.
	s := grpc.NewServer()

	// Register the LogServiceServer with the gRPC server.
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	// Serve incoming connections.
	if err := s.Serve(lis); err != nil {
		return err
	}
	return nil
}
