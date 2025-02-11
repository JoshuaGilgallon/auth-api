package handlers

import (
	"github.com/gin-gonic/gin"
	"auth-api/internal/services"
	"net/http"
)

func GetUsers(c *gin.Context) {
	users := services.FetchUsers()
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var userInput services.UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := services.CreateUser(userInput)
	c.JSON(http.StatusCreated, user)
}
