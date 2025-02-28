package handlers

import (
	"auth-api/internal/errors"
	"auth-api/internal/services"
	"auth-api/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body services.UserInput true "User input data"
// @Success 201 {object} string "user created"
// @Router /api/user [post]
func CreateUser(c *gin.Context) {
	var userInput services.UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Sanitize inputs
	userInput.Email = strings.TrimSpace(userInput.Email)
	userInput.PhoneNumber = strings.TrimSpace(userInput.PhoneNumber)

	// Validate required fields
	if userInput.Email == "" && userInput.PhoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or phone number is required"})
		return
	}

	if !utils.IsValidPassword(userInput.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password does not meet security requirements"})
		return
	}

	user, err := services.CreateUser(userInput)
	if err != nil {

		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.AlreadyExists:
				c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			case errors.ValidationError:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

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
