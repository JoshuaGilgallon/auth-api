package handlers

import (
	"auth-api/internal/errors"
	"auth-api/internal/services"
	"auth-api/internal/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Login
// @Description Log into a user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body services.LoginInput false "Login Input"
// @Success 200 {string} string "successfully logged in"
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	var loginInput services.LoginInput

	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Sanitize inputs
	loginInput.Email = strings.TrimSpace(loginInput.Email)
	loginInput.PhoneNumber = strings.TrimSpace(loginInput.PhoneNumber)

	// Validate that at least one login method is provided
	if loginInput.Email == "" && loginInput.PhoneNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or phone number is required"})
		return
	}

	session, err := services.Login(loginInput)
	if err != nil {
		log.Printf("Login attempt failed: %v", err)

		switch e := err.(type) {
		case *errors.LoginError:
			switch e.Type {
			case errors.InvalidCredentials:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			case errors.AccountLocked:
				c.JSON(http.StatusForbidden, gin.H{"error": "Account is locked"})
			case errors.AccountDisabled:
				c.JSON(http.StatusForbidden, gin.H{"error": "Account is disabled"})
			case errors.TooManyAttempts:
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many login attempts"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed"})
			}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// Use email or phone number for logging, whichever was provided
	identifier := loginInput.Email
	if identifier == "" {
		identifier = loginInput.PhoneNumber
	}
	log.Printf("Successful login for user with identifier: %s", identifier)
	c.JSON(http.StatusOK, session)
}

// @Summary Logout
// @Description Logout of a user account
// @Tags auth
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer <token>"
// @Success 200 {string} string "successfully logged out"
// @Router /api/auth/logout [post]
func Logout(c *gin.Context) {
	token := utils.ExtractBearerToken(c.GetHeader("Authorization"))
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := services.Logout(token); err != nil {
		log.Printf("Logout error: %v", err)

		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.SessionNotFound:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired session"})
			case errors.InvalidToken:
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token format"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
			}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// @Summary Sign up
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body services.UserInput true "User input data"
// @Success 201 {object} models.User
// @Router /api/auth/signup [post]
func SignUp(c *gin.Context) {
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
		log.Printf("User creation error: %v", err)

		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.AlreadyExists:
				c.JSON(http.StatusConflict, gin.H{"error": "Email or phone number already exists"})
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

	// Update logging to use email or phone number instead of username
	identifier := userInput.Email
	if identifier == "" {
		identifier = userInput.PhoneNumber
	}
	log.Printf("New user created with identifier: %s", identifier)
	c.JSON(http.StatusCreated, user)
}
