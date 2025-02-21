package handlers

import (
	"auth-api/internal/errors"
	"auth-api/internal/models"
	"auth-api/internal/services"
	"net/http"
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
	email := c.Query("email")
	phoneNumber := c.Query("phone_number")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	updatedStartTime := c.Query("updated_start_time")
	updatedEndTime := c.Query("updated_end_time")

	var users []models.User

	if startTime != "" && endTime != "" {
		start, err := time.Parse(time.RFC3339, startTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format"})
			return
		}
		end, err := time.Parse(time.RFC3339, endTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format"})
			return
		}
		users, err = services.SearchUserByCreateTimeRange(start, end)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if updatedStartTime != "" && updatedEndTime != "" {
		start, err := time.Parse(time.RFC3339, updatedStartTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid updated start time format"})
			return
		}
		end, err := time.Parse(time.RFC3339, updatedEndTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid updated end time format"})
			return
		}
		users, err = services.SearchUsersByTimeUpdatedRange(start, end)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if email != "" {
		user, err := services.SearchUserByCredentials(email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	if phoneNumber != "" {
		user, err := services.SearchUserByCredentials(phoneNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	// remove duplicate users
	userMap := make(map[string]models.User)
	for _, user := range users {
		userMap[user.ID.Hex()] = user
	}

	uniqueUsers := make([]models.User, 0, len(userMap))
	for _, user := range userMap {
		uniqueUsers = append(uniqueUsers, user)
	}

	c.JSON(http.StatusOK, uniqueUsers)
}

// when creating an administrator session,
// ensure that there is no token stored in a cookie.
// These tokens refresh after every use due to their administrator nature
// and therefore need to be more secure. You will be logged out after 30 minutes
// of inactivity.

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

	session, err := services.AdminLogin(adminLoginInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization token"})
		return
	}

	// Remove "Bearer " prefix if present
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}

	if err := services.Logout(accessToken); err != nil {
		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.SessionNotFound:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired session"})
			case errors.InvalidToken:
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token format"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
			}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}
