package server

import (
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Logging middleware
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		log.Printf("Request processed in %v\n", time.Since(startTime))
	}
}

var (
	ipLimiter = make(map[string]*rate.Limiter)
	mu        sync.RWMutex
)

func getRateLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := ipLimiter[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(1*time.Second), 3)
		ipLimiter[ip] = limiter
	}
	return limiter
}

func SecurityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Content-Security-Policy", "default-src 'self'")

		// Rate limiting
		ip := c.ClientIP()
		limiter := getRateLimiter(ip)
		if !limiter.Allow() {
			c.AbortWithStatus(429)
			return
		}

		c.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		c.Next()
	}
}
