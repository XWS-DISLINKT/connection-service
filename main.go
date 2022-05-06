package main

import (
	"connection-service/startup"
	"connection-service/startup/config"
)

func main() {
	config := config.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
