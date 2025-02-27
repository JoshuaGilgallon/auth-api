package main

import (
	internalconfig "auth-api/internal/config"
	"auth-api/internal/repositories"
	"auth-api/internal/services"
	"log"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var (
	nextPurgeTime time.Time
	mu            sync.Mutex
)

func main() {
	log.Println("Worker service starting...")

	log.Println("Connecting to database...")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConfig := internalconfig.NewDatabaseConfig()
	if err := repositories.InitDatabase(dbConfig); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repositories.CloseDatabase()

	log.Println("Successfully connected to database")

	log.Println("Worker service successfully started!")
	go func() {
		for {
			mu.Lock()
			nextPurgeTime = time.Now().Add(24 * time.Hour)
			services.UpdatePurgeTime(nextPurgeTime)
			mu.Unlock()

			err := services.DeleteInvalidSessions()
			if err != nil {
				log.Println("WORKER SERVICE: Purging old sessions...", err)
			} else {
				log.Println("WORKER SERVICE: Old sessions purged successfully")
			}

			time.Sleep(24 * time.Hour)
		}
	}()

	// keep worker alive
	select {}
}
