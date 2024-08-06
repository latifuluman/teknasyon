package grpc

import (
	"google.golang.org/grpc"

	"listener/grpc/logs"
	"listener/grpc/mail"
)

var (
	mailConn   mail.MailServiceClient //mail-service gRPC connection client
	loggerConn logs.LogServiceClient  //logger-service gRPC connection client
)

func SetGrpcMailConn(conn *grpc.ClientConn) {
	mailConn = mail.NewMailServiceClient(conn)
}
func SetGrpcLoggerConn(conn *grpc.ClientConn) {
	loggerConn = logs.NewLogServiceClient(conn)
}

func GetGrpcMailClient() mail.MailServiceClient {
	return mailConn
}

func GetGrpcLoggerClient() logs.LogServiceClient {
	return loggerConn
}
