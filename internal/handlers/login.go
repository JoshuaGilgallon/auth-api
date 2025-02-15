package handlers

import (
	"auth-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Login
// @Description Log into a user account
// @Tags login users
// @Accept json
// @Produce json
// @Param user body services.LoginInput false "Login Input"
// @Success 200 {string} string "successfully logged in"
// @Router /api/login [post]
func Login(c *gin.Context) {
	var loginInput services.LoginInput

	// Bind JSON request body to loginInput struct
	if err := c.ShouldBindJSON(&loginInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	session, err := services.Login(loginInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, session)
}
