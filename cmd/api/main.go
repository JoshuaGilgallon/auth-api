package main

import (
	"log"
	"auth-api/config"
	"auth-api/internal/server"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize and start the API server
	r := server.SetupRouter()
	log.Println("Starting server on port 8080...")
	r.Run(":8080")
}
