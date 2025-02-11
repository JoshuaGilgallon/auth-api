// cmd/api/main.go
package main

import (
    "log"
    "auth-api/config"
    "auth-api/internal/server"
    "auth-api/docs"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Auth API
// @version 1.0
// @description Authentication and Authorization API
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
    // Load configuration
    config.LoadConfig()

    docs.SwaggerInfo.Title = "Auth API"
    docs.SwaggerInfo.Description = "Authentication and Authorization API"
    docs.SwaggerInfo.Version = "1.0"
    docs.SwaggerInfo.Host = "localhost:8080"
    docs.SwaggerInfo.BasePath = "/"
    docs.SwaggerInfo.Schemes = []string{"http", "https"}

    // Initialize and start the API server
    r := server.SetupRouter()

    // Add swagger endpoint
    r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    log.Println("Starting server on port 8080...")
    r.Run(":8080")
}