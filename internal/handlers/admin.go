package handlers

import (
	"auth-api/internal/errors"
	"auth-api/internal/models"
	"auth-api/internal/services"
	"auth-api/internal/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Advanced search for users
// @Description Search for users using various criteria
// @Tags admin
// @Accept json
// @Produce json
// @Param email query string false "Email"
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

	// Build search criteria
	criteria := models.UserAdvancedSearchCriteria{
		FirstName: strings.TrimSpace(c.Query("first_name")),
		LastName:  strings.TrimSpace(c.Query("last_name")),
		Email:     strings.TrimSpace(c.Query("email")),
	}

	// Parse time parameters
	if startTime := c.Query("start_time"); startTime != "" {
		if parsed, err := utils.SafeParseTime(startTime); err == nil {
			criteria.StartTime = &parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
			return
		}
	}

	if endTime := c.Query("end_time"); endTime != "" {
		if parsed, err := utils.SafeParseTime(endTime); err == nil {
			criteria.EndTime = &parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
			return
		}
	}

	if updateStartTime := c.Query("updated_start_time"); updateStartTime != "" {
		if parsed, err := utils.SafeParseTime(updateStartTime); err == nil {
			criteria.UpdateStartTime = &parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid updated_start_time format"})
			return
		}
	}

	if updateEndTime := c.Query("updated_end_time"); updateEndTime != "" {
		if parsed, err := utils.SafeParseTime(updateEndTime); err == nil {
			criteria.UpdateEndTime = &parsed
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid updated_end_time format"})
			return
		}
	}

	// Perform search
	users, err := services.SearchUsers(criteria)
	if err != nil {
		log.Printf("Search error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search operation failed"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Simple search for users
// @Description Search for users using a single search term across all fields
// @Tags admin
// @Accept json
// @Produce json
// @Param search query string true "Search term"
// @Success 200 {array} models.UserSearchCriteria
// @Router /api/admin/search [get]
func SimpleSearch(c *gin.Context) {
	// Validate admin session first
	adminSession, err := validateAdminSession(c)
	if err != nil {
		return
	}

	var userSearchCriteria models.UserSearchCriteria

	// Bind JSON request body to UserSearchCriteria struct
	if err := c.ShouldBindQuery(&userSearchCriteria); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format:"})
		return
	}

	// Log admin action
	log.Printf("Admin %s performing simple search", adminSession.AdminID)

	// Get search term from query
	searchTerm := strings.TrimSpace(userSearchCriteria.SearchTerm)
	if searchTerm == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search term is required"})
		return
	}

	// Perform search
	result, err := services.SimpleSearch(searchTerm, userSearchCriteria.PageNumber, userSearchCriteria.PageSize)
	if err != nil {
		log.Printf("Search error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Search operation failed"})
		return
	}

	// Return both users and total results
	c.JSON(http.StatusOK, gin.H{
		"users":         result.Users,
		"total_results": result.TotalResults,
	})
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Set secure HTTP-only cookie
	c.SetCookie(
		"admin_token",
		session.(map[string]interface{})["access_token"].(string),
		1800, // 30 minutes
		"/",
		"",
		true, // secure
		true, // httpOnly
	)

	// Don't send the token in the response body
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// @Summary Admin Logout
// @Description Logout of an admin account
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {string} string "successfully logged out"
// @Router /api/admin/logout [post]
func AdminLogout(c *gin.Context) {
	// Clear the cookie
	c.SetCookie(
		"admin_token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

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
// @Success 200 {object} models.AdminSession
// @Router /api/admin/validate [get]
func ValidateAdminSession(c *gin.Context) {
	token, err := c.Cookie("admin_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
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
	token, err := c.Cookie("admin_token")
	if err != nil {
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
