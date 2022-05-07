package startup

import (
	"connection-service/application"
	"connection-service/infrastructure/api"
	"connection-service/startup/config"
	"fmt"
	connections "github.com/XWS-DISLINKT/dislinkt/common/proto/connection-service"
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

func (server *Server) initHandlers() {

}

func (server *Server) Start() {
	connectionService := application.NewConnectionService()
	connectionHandler := server.initConnectionsHandler(connectionService)
	connectionHandler.Demo()
	server.startGrpcServer(connectionHandler)
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
