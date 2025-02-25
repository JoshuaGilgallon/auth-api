package handlers

import (
	"auth-api/internal/errors"
	"auth-api/internal/models"
	"auth-api/internal/services"
	"auth-api/internal/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Advanced search for users
// @Description Search for users using various criteria
// @Tags admin
// @Accept json
// @Produce json
// @Param email query string false "Email"
// @Param phone_number query string false "Phone Number"
// @Param start_time query string false "Start Time" format(date-time)
// @Param end_time query string false "End Time" format(date-time)
// @Param updated_start_time query string false "Updated Start Time" format(date-time)
// @Param updated_end_time query string false "Updated End Time" format(date-time)
// @Success 200 {array} models.User
// @Router /api/admin/advancedSearch [get]
func AdvancedSearch(c *gin.Context) {
	// Validate admin session first
	adminSession, err := validateAdminSession(c)
	if err != nil {
		return
	}

	// Log admin action
	log.Printf("Admin %s performing advanced search", adminSession.AdminID)

	// Sanitize inputs
	email := strings.TrimSpace(c.Query("email"))
	phoneNumber := strings.TrimSpace(c.Query("phone_number"))

	// Validate time formats
	timeParams := map[string]string{
		"start_time":         c.Query("start_time"),
		"end_time":           c.Query("end_time"),
		"updated_start_time": c.Query("updated_start_time"),
		"updated_end_time":   c.Query("updated_end_time"),
	}

	times := make(map[string]time.Time)
	for key, value := range timeParams {
		if value != "" {
			parsedTime, err := utils.SafeParseTime(value)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format"})
				return
			}
			times[key] = parsedTime
		}
	}

	var users []models.User
	var searchErr error

	// Perform searches based on available parameters
	if !times["start_time"].IsZero() && !times["end_time"].IsZero() {
		users, searchErr = services.SearchUserByCreateTimeRange(times["start_time"], times["end_time"])
	} else if !times["updated_start_time"].IsZero() && !times["updated_end_time"].IsZero() {
		users, searchErr = services.SearchUsersByTimeUpdatedRange(times["updated_start_time"], times["updated_end_time"])
	} else if email != "" {
		var user models.User
		user, searchErr = services.SearchUserByCredentials(email)
		if searchErr == nil {
			users = append(users, user)
		}
	} else if phoneNumber != "" {
		var user models.User
		user, searchErr = services.SearchUserByCredentials(phoneNumber)
		if searchErr == nil {
			users = append(users, user)
		}
	}

	if searchErr != nil {
		log.Printf("Search error: %v", searchErr)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search operation failed"})
		return
	}

	// Deduplicate users
	uniqueUsers := utils.DeduplicateUsers(users)
	c.JSON(http.StatusOK, uniqueUsers)
}

// when creating an administrator session, ensure that there is no token stored in a cookie.
// These tokens refresh after every use due to their administrator nature and therefore need to be more secure. 
// You will be logged out after 30 minutes of inactivity.

// @Summary Administrator portal login
// @Description Allows an administrator to login to the admin portal
// @Tags admin
// @Accept json
// @Produce json
// @Param user body services.AdminLoginInput false "Admin Login Input"
// @Success 200 {string} string "Logged in"
// @Router /api/admin/login [post]
func AdminLogin(c *gin.Context) {
	var adminLoginInput services.AdminLoginInput

	// Bind JSON request body to loginInput struct
	if err := c.ShouldBindJSON(&adminLoginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Sanitize inputs
	adminLoginInput.Username = strings.TrimSpace(adminLoginInput.Username)
	if !utils.IsValidUsername(adminLoginInput.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username format"})
		return
	}

	session, err := services.AdminLogin(adminLoginInput)
	if err != nil {
		// Don't expose specific error messages
		log.Printf("Login error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Log successful login
	log.Printf("Admin login successful: %s", adminLoginInput.Username)
	c.JSON(http.StatusOK, session)
}

// @Summary Admin Logout
// @Description Logout of an admin account
// @Tags admin
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer <token>"
// @Success 200 {string} string "successfully logged out"
// @Router /api/admin/logout [post]
func AdminLogout(c *gin.Context) {
	accessToken := utils.ExtractBearerToken(c.GetHeader("Authorization"))
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := services.AdminLogout(accessToken); err != nil {
		log.Printf("Logout error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// Root User can only do the following endpoints

// @Summary Create admin user
// @Description Create a new admin account - Only for root user
// @Tags admin
// @Accept json
// @Produce json
// @Param request body services.AdminCreationRequest true "Admin creation request with admin user and root credentials"
// @Success 201 {object} models.AdminUser
// @Router /api/admin/create [post]
func CreateAdminAccount(c *gin.Context) {
	var request services.AdminCreationRequest

	// Bind the JSON request body
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Sanitize and validate inputs
	request.AdminUser.Username = strings.TrimSpace(request.AdminUser.Username)
	request.RootUser.Username = strings.TrimSpace(request.RootUser.Username)

	if !utils.IsValidUsername(request.AdminUser.Username) ||
		!utils.IsValidPassword(request.AdminUser.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username or password format"})
		return
	}

	// Validate the root user credentials
	if !services.ValidateRootUserCredentials(request.RootUser.Username, request.RootUser.Password) {
		// Don't expose specific error messages
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Create the admin user
	user, err := services.CreateAdminUser(request.AdminUser)
	if err != nil {
		log.Printf("Admin creation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin user"})
		return
	}

	// Log admin creation
	log.Printf("New admin created by root user: %s", request.AdminUser.Username)
	c.JSON(http.StatusCreated, user)
}

// @Summary Validate admin access token
// @Description Validate an access token and return session info
// @Tags admin
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer <token>"
// @Success 200 {object} models.AdminSession
// @Router /api/admin/validate [get]
func ValidateAdminSession(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.GetHeader("Token")
	}

	adminSession, err := services.ValidateAdminAccessToken(token)
	if err != nil {
		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.SessionNotFound:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired session"})
			case errors.InvalidToken:
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token format"})
			case errors.TokenExpired:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate session"})
			}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, adminSession)
}

func validateAdminSession(c *gin.Context) (models.AdminSession, error) {
	token := utils.ExtractBearerToken(c.GetHeader("Authorization"))
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return models.AdminSession{}, errors.NewAuthenticationError("missing token", nil)
	}

	session, err := services.ValidateAdminAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return models.AdminSession{}, err
	}

	return session, nil
}
