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

// @Summary Get user by ID
// @Description Get a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Router /api/user/{id} [get]
func GetUser(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	user, err := services.GetUser(id)
	if err != nil {

		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			case errors.ValidationError:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	if user.ID.IsZero() {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Get user info via access token
// @Description Returns information about the curently authenticated user via their access token
// @Tags users
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer <token>"
// @Success 200 {object} models.User
// @Router /api/user/me [get]
func GetCurrentUser(c *gin.Context) {
	token := utils.ExtractBearerToken(c.GetHeader("Authorization"))
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization token provided"})
		return
	}

	user, err := services.GetCurrentUser(token)
	if err != nil {

		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.SessionNotFound, errors.TokenExpired:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			case errors.InvalidToken:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token format"})
			case errors.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user information"})
			}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// ADMIN USER ENDPOINTS:

// @Summary Edit user
// @Description Edits user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UserInput true "User input data"
// @Success 200 {object} models.User
// @Router /api/admin/updateuser [patch]
func AdminUpdateUser(c *gin.Context) {
	adminSession, err := validateAdminSession(c)
	if err != nil {
		log.Printf("Admin session validation failed: %v", err)
		return
	}

	log.Printf("Admin %s updating user", adminSession.AdminID)

	userID := c.Query("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	var userInput models.FullUserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		log.Printf("Error binding JSON: %v, Body: %v", err, c.Request.Body)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	log.Printf("Received update request for user %s with data: %+v", userID, userInput)

	updatedUser, err := services.UpdateUser(userID, userInput)
	if err != nil {
		log.Printf("Error updating user %s: %v", userID, err)
		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.NotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			case errors.ValidationError:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": e.Error()})
			}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
