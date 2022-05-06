package startup

import (
	"connection-service/application"
	"connection-service/infrastructure/api"
	"connection-service/startup/config"
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
	connectionHandler := api.NewConnectionHandler(connectionService)
	connectionHandler.Demo()
}
