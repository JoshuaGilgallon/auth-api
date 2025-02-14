package internalconfig

import (
	"os"
	"time"
)

type DatabaseConfig struct {
	URI             string
	ConnectTimeout  time.Duration
	DatabaseName    string
}

func NewDatabaseConfig() *DatabaseConfig {
	uri := os.Getenv("DATABASE_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	return &DatabaseConfig{
		URI:            uri,
		ConnectTimeout: 30 * time.Second,  // Increased timeout
		DatabaseName:   "auth_api",
	}
}