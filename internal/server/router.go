package server

import (
	"github.com/gin-gonic/gin"
	"auth-api/internal/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", handlers.HealthCheck)

	// Grouped API endpoints
	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("/", handlers.GetUsers)
			users.POST("/", handlers.CreateUser)
		}
	}

	return r
}