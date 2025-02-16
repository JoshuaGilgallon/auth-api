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
			users.GET("/:id", handlers.GetUser)
		}
		sessions := api.Group("/session")
		{
			sessions.POST("/", handlers.CreateSession)
			sessions.GET("/validate", handlers.ValidateSession)
			sessions.POST("/refresh", handlers.RefreshSession)
			sessions.DELETE("/:session_id", handlers.InvalidateSession)
		}

		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", handlers.Logout)
			auth.POST("/signup", handlers.SignUp)
		}
	}

	return r
}
