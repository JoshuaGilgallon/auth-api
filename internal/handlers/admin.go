package handlers

import (
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
// @Param username query string false "Email"
// @Param password query string false "Phone Number"
// @Success 200 {string} "Logged in"
// @Router /api/admin/login [post]
func AdminLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Logged in"}) // temporary, add logic later (im in class rn lol)
}
