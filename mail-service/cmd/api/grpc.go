package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	mail "mail-service/grpc"
)

type MailServer struct {
	mail.UnimplementedMailServiceServer
	Config *Config
}

// SendMail receives message info via grpc and sends mail to specified user
func (m *MailServer) SendMail(ctx context.Context, req *mail.MailRequest) (*mail.MailResponse, error) {
	from := req.GetFrom()
	to := req.GetTo()
	subject := req.GetSubject()
	message := req.GetMessage()

	msg := Message{
		From:    from,
		To:      to,
		Subject: subject,
		Message: message,
	}

	err := m.Config.Mailer.SendSMTPMessage(msg)
	if err != nil {
		return &mail.MailResponse{}, nil
	}

	payload := &mail.MailResponse{

		Result: "sent to " + msg.To,
	}
	return payload, nil
}

// gRPCListen starts the gRPC server and listens for incoming connections.
func (app *Config) gRPCListen() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	mail.RegisterMailServiceServer(s, &MailServer{Config: app})

	if err := s.Serve(lis); err != nil {

		return err
	}
	return nil
}
