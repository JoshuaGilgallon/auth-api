package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Health check
// @Description Checks health of server
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
