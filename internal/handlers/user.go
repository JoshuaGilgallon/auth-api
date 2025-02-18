package handlers

import (
	"auth-api/internal/services"
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.CreateUser(userInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	id := c.Param("id")
	user, err := services.GetUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Param token header string true "Access Token"
// @Success 200 {object} models.User
// @Router /api/user/me [get]
func GetCurrentUser(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.GetHeader("Token")
	}

	user, err := services.GetCurrentUser(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}