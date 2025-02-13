package internalconfig

import "time"

type DatabaseConfig struct {
	URI             string
	ConnectTimeout  time.Duration
	DatabaseName    string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		URI:            "mongodb://localhost:27017",
		ConnectTimeout: 10 * time.Second,
		DatabaseName:   "auth_api",
	}
}