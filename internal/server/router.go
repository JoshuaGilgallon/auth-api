package server

import (
	"auth-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", handlers.HealthCheck)

	// Grouped API endpoints
	api := r.Group("/api")
	{
		users := api.Group("/user")
		{
			users.POST("/", handlers.CreateUser)
			users.GET("/:id", handlers.GetUser) // New endpoint
		}
		sessions := api.Group("/session")
		{
			sessions.POST("/", handlers.CreateSession)
			sessions.GET("/validate", handlers.ValidateSession)
			sessions.POST("/refresh", handlers.RefreshSession)
			sessions.DELETE("/:session_id", handlers.InvalidateSession)
		}
		login := api.Group("/login")
		{
			login.POST("/", handlers.Login)
		}
	}

	return r
}
