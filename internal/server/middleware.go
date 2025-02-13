package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// Logging middleware
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		fmt.Printf("Request processed in %v\n", time.Since(startTime))
	}
}