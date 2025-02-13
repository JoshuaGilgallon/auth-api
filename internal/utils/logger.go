package utils

import (
	"log"
)

func LogInfo(message string) {
	log.Println("[INFO]", message)
}

func LogError(message string, err error) {
	log.Println("[ERROR]", message, err)
}