package startup

import (
	"connection-service/application"
	"connection-service/infrastructure/api"
	"connection-service/startup/config"
	"fmt"
	connections "github.com/XWS-DISLINKT/dislinkt/common/proto/connection-service"
	saga "github.com/XWS-DISLINKT/dislinkt/common/saga/messaging"
	"github.com/XWS-DISLINKT/dislinkt/common/saga/messaging/nats"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	server := &Server{
		config: config,
	}
	server.initHandlers()
	return server
}

const (
	QueueGroup = "connection_service"
)

func (server *Server) initHandlers() {

}

func (server *Server) Start() {
	connectionService := application.NewConnectionService()

	commandSubscriber := server.initSubscriber(server.config.CreateUserCommandSubject, QueueGroup)
	replyPublisher := server.initPublisher(server.config.CreateUserReplySubject)
	server.initCreateOrderHandler(connectionService, replyPublisher, commandSubscriber)

	connectionHandler := server.initConnectionsHandler(connectionService)
	connectionHandler.Demo()
	server.startGrpcServer(connectionHandler)
}

func (server *Server) initCreateOrderHandler(service *application.ConnectionService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewCreateUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initPublisher(subject string) saga.Publisher {
	publisher, err := nats.NewNATSPublisher(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}

func (server *Server) initSubscriber(subject, queueGroup string) saga.Subscriber {
	subscriber, err := nats.NewNATSSubscriber(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject, queueGroup)
	if err != nil {
		log.Fatal(err)
	}
	return subscriber
}

func (server *Server) initConnectionsHandler(service *application.ConnectionService) *api.ConnectionHandler {
	return api.NewConnectionHandler(service)
}

func (server *Server) startGrpcServer(connectionHandler *api.ConnectionHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	connections.RegisterConnectionServiceServer(grpcServer, connectionHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
