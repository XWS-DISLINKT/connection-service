package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port          string
	Neo4jPort     string
	Neo4jHost     string
	Neo4jProtocol string
	Neo4jUsername string
	Neo4jPassword string
}

func NewConfig() *Config {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		fmt.Println("docker")

		return &Config{
			Port:          os.Getenv("CONNECTION_SERVICE_PORT"),
			Neo4jPort:     os.Getenv("CONNECTION_DB_PORT"),
			Neo4jHost:     os.Getenv("CONNECTION_DB_HOST"),
			Neo4jProtocol: os.Getenv("CONNECTION_DB_PROTOCOL"),
			Neo4jUsername: os.Getenv("CONNECTION_DB_USERNAME"),
			Neo4jPassword: os.Getenv("CONNECTION_DB_PASSWORD"),
		}
	} else {
		fmt.Println("local")

		return &Config{
			Port:          "8004",
			Neo4jPort:     "7687",
			Neo4jHost:     "localhost",
			Neo4jProtocol: "bolt",
			Neo4jUsername: "neo4j",
			Neo4jPassword: "password",
		}
	}
}
