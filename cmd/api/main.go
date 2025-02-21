// cmd/api/main.go
package main

import (
	"auth-api/docs"
	internalconfig "auth-api/internal/config"
	"auth-api/internal/repositories"
	"auth-api/internal/server"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Auth API
// @version 1.2
// @description Authentication and Authorization API
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Swagger docs
	docs.SwaggerInfo.Title = "Auth API"
	docs.SwaggerInfo.Description = "Authentication and Authorization API"
	docs.SwaggerInfo.Version = "Beta 1.2"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Initialize the database
	dbConfig := internalconfig.NewDatabaseConfig()
	if err := repositories.InitDatabase(dbConfig); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repositories.CloseDatabase()

	// Initialize and start the API server
	r := server.SetupRouter()
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Define routes for all admin pages
	r.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{"title": "Admin Dashboard"})
	})

	r.GET("/admin/users", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users.html", gin.H{"title": "User Management"})
	})

	r.GET("/admin/sessions", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sessions.html", gin.H{"title": "Active Sessions"})
	})

	r.GET("/admin/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_login.html", gin.H{"title": "Admin Login"})
	})

	r.SetTrustedProxies(nil)

	log.Println("Starting server on port 8080...")
	r.Run(":8080")
}
