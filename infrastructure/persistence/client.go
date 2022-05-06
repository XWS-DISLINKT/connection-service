package persistence

import (
	"connection-service/startup/config"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func GetDriver() (neo4j.Driver, error) {
	cfg := config.NewConfig()
	uri := fmt.Sprintf("%s://%s:%s/", cfg.Neo4jProtocol, cfg.Neo4jHost, cfg.Neo4jPort)
	return neo4j.NewDriver(uri, neo4j.BasicAuth(cfg.Neo4jUsername, cfg.Neo4jPassword, ""))
}
