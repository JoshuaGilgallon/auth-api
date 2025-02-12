package handlers

import (
	"github.com/gin-gonic/gin"
	"auth-api/internal/services"
	"net/http"
)


// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} string "user created"
// @Router /api/user [post]
func CreateUser(c *gin.Context) {
	var userInput services.UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := services.CreateUser(userInput)
	c.JSON(http.StatusCreated, user)
}
