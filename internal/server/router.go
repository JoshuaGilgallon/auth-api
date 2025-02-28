package server

import (
	"auth-api/internal/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Serve static files
	r.Static("/internal/static", "./internal/static")

	// Load templates in the correct order - base templates first, then pages
	r.LoadHTMLFiles(
		"templates/partials/sidebar.html",
		"templates/dashboard-content.html",
		"templates/users.html",
		"templates/sessions.html",
		"templates/admin_login.html",
		"templates/auth-check.html",
	)

	// Enable debug mode to see more detailed logs
	gin.SetMode(gin.DebugMode)

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

		admin := api.Group("/admin")
		{
			admin.GET("/advancedSearch", handlers.AdvancedSearch)
			admin.POST("/login", handlers.AdminLogin)
			admin.POST("/logout", handlers.AdminLogout)
			admin.POST("/create", handlers.CreateAdminAccount)
			admin.GET("/validate", handlers.ValidateAdminSession)
		}
		stats := api.Group("/stats")
		{
			stats.GET("/cache", handlers.GetCacheStats)
			stats.GET("/dashboard", handlers.GetDashboardStats)
		}
	}

	// Admin web routes
	admin := r.Group("/admin")
	{
		admin.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "admin_login.html", gin.H{"title": "Admin Login"})
		})
		admin.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "auth-check.html", gin.H{})
		})
		admin.GET("/dashboard/content", func(c *gin.Context) {
			c.HTML(http.StatusOK, "dashboard-content.html", gin.H{
				"title":  "Dashboard",
				"active": "dashboard", // This will highlight the current page in sidebar
			})
		})
		admin.GET("/users", func(c *gin.Context) {
			c.HTML(http.StatusOK, "users.html", gin.H{"title": "User Management"})
		})
		admin.GET("/sessions", func(c *gin.Context) {
			c.HTML(http.StatusOK, "sessions.html", gin.H{"title": "Active Sessions"})
		})
	}

	return r
}
