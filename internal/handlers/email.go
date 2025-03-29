package handlers

import (
	"auth-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Email redirect
// @Description Gets the code from the URL path and resumes the signup process
// @Tags email
// @Accept json
// @Produce json
// @Param code path string true "Verification code"
// @Success 200
// @Router /api/email/verify/{code} [get]
func EmailRedirect(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification code is required"})
		return
	}

	url, err := services.VerifyEmail(code)
	if err != nil || url == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	c.Redirect(http.StatusFound, url)
}

// @Summary Email verify NO redirect
// @Description Verifies the code without redirecting
// @Tags email
// @Accept json
// @Produce json
// @Param code path string true "Verification code"
// @Success 200
// @Router /api/email//verify/NR/{code} [get]
func EmailVerifyNoRedirect(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification code is required"})
		return
	}

	url, err := services.VerifyEmail(code)
	if err != nil || url == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
