package handlers

import (
	"auth-api/internal/errors"
	"auth-api/internal/models"
	"auth-api/internal/repositories"
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
// @Param user body models.LoginInput false "Login Input"
// @Success 200
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	var loginInput models.LoginInput

	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Sanitize inputs
	loginInput.Email = strings.TrimSpace(loginInput.Email)

	if loginInput.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}

	session, err := services.Login(loginInput)
	if err != nil {
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

	repositories.IncreaseLoginCount()

	c.SetCookie("refresh_token", session.RefreshToken, int(session.RefreshExpiresAt.Unix()), "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Successfully logged in!",
		"access_token": session.AccessToken,
		"expires_at":   session.AccessExpiresAt,
	})
}

// @Summary Logout
// @Description Logout of a user account
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {string} string "successfully logged out"
// @Router /api/auth/logout [post]
func Logout(c *gin.Context) {
	// clear the cookies by setting them to expire immediately
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	token, err := c.Cookie("access_token")
	if err != nil && err != http.ErrNoCookie {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read cookie"})
		return
	}

	// only attempt logout if we have a token
	if token != "" {
		if err := services.Logout(token); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate session"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// @Summary Sign up
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserInput true "User input data"
// @Success 201
// @Router /api/auth/signup [post]
func SignUp(c *gin.Context) {
	var userInput models.UserInput

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// sanitise inputs
	userInput.Email = strings.TrimSpace(userInput.Email)

	// validate required fields
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
				c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			case errors.ValidationError:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user, please try again later"})
			}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error, please try again later"})
		}
		return
	}

	log.Printf("User created: %v", user)
	log.Printf("User ID: %s", user.ID.Hex())

	verifEmail := models.VerifEmailInput{
		UserID: user.ID.Hex(),
	}

	response, err := services.CreateVerifEmail(verifEmail)
	if err != nil {
		log.Printf("Error sending verification email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
		return
	}
	log.Printf("Response %v", response)

	user_return := models.UserCreateReturn{
		Success: true,
		Message: "User created successfully",
	}

	c.JSON(http.StatusCreated, user_return)
}

// @Summary Complete Sign up
// @Description Finishes off the signup process after email verification succeeds + logs in the user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.SetupUserInput true "Setup step 2 user input data"
// @Success 201
// @Router /api/auth/csignup [post]
func FinishSignup(c *gin.Context) {
	var userInput models.SetupUserInput

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// sanitise inputs
	userInput.FirstName = strings.TrimSpace(userInput.FirstName)
	userInput.LastName = strings.TrimSpace(userInput.LastName)

	// validate required fields
	if userInput.FirstName == "" || userInput.LastName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "First name and last name are required"})
		return
	}

	// check birthdate is in the YYY-MM-DD format
	if !utils.IsValidDate(userInput.BirthDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid birthdate format"})
		return
	}

	response, err := services.CompleteSignup(userInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete signup"})
		log.Printf("Error completing signup: %v", err)
		return
	}

	// delete the token from the database so it cant be used again
	err = repositories.InvalidateEmailToken(userInput.Token)
	if err != nil {
		log.Printf("Error invalidating email token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate email token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Signup completed successfully", "access": response.AccessToken})

	c.SetCookie("refresh_token", response.RefreshToken, int(response.RefreshExpiresAt.Unix()), "/", "", true, true)
}
