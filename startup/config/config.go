package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port                     string
	Neo4jPort                string
	Neo4jHost                string
	Neo4jProtocol            string
	Neo4jUsername            string
	Neo4jPassword            string
	ProfileServiceHost       string
	ProfileServicePort       string
	NatsHost                 string
	NatsPort                 string
	NatsUser                 string
	NatsPass                 string
	CreateUserCommandSubject string
	CreateUserReplySubject   string
}

func NewConfig() *Config {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		fmt.Println("docker")

		return &Config{
			Port:                     os.Getenv("CONNECTION_SERVICE_PORT"),
			Neo4jPort:                os.Getenv("CONNECTION_DB_PORT"),
			Neo4jHost:                os.Getenv("CONNECTION_DB_HOST"),
			Neo4jProtocol:            os.Getenv("CONNECTION_DB_PROTOCOL"),
			Neo4jUsername:            os.Getenv("CONNECTION_DB_USERNAME"),
			Neo4jPassword:            os.Getenv("CONNECTION_DB_PASSWORD"),
			ProfileServiceHost:       os.Getenv("PROFILE_SERVICE_HOST"),
			ProfileServicePort:       os.Getenv("PROFILE_SERVICE_PORT"),
			NatsHost:                 os.Getenv("NATS_HOST"),
			NatsPort:                 os.Getenv("NATS_PORT"),
			NatsUser:                 os.Getenv("NATS_USER"),
			NatsPass:                 os.Getenv("NATS_PASS"),
			CreateUserCommandSubject: os.Getenv("CREATE_USER_COMMAND_SUBJECT"),
			CreateUserReplySubject:   os.Getenv("CREATE_USER_REPLY_SUBJECT"),
		}
	} else {
		fmt.Println("local")

		return &Config{
			Port:                     "8004",
			Neo4jPort:                "7687",
			Neo4jHost:                "localhost",
			Neo4jProtocol:            "bolt",
			Neo4jUsername:            "neo4j",
			Neo4jPassword:            "password",
			ProfileServiceHost:       "localhost",
			ProfileServicePort:       "8001",
			NatsHost:                 "localhost",
			NatsPort:                 "4222",
			NatsUser:                 "ruser",
			NatsPass:                 "T0pS3cr3t",
			CreateUserCommandSubject: "user.create.command",
			CreateUserReplySubject:   "user.create.reply",
		}
	}
}
