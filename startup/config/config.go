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
			Port:          "8000",
			Neo4jPort:     "7687",
			Neo4jHost:     "connection-db",
			Neo4jProtocol: "bolt",
			Neo4jUsername: "neo4j",
			Neo4jPassword: "neo4j",
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
