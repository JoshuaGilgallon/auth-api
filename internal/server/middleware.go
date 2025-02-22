package server

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logging middleware
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		log.Printf("Request processed in %v\n", time.Since(startTime))
	}
}
