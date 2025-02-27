package handlers

import (
	"auth-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// All these endpoints need to be called with an admin access token:

// @Summary Get Dashboard Stats
// @Description Gets information to display on the admin dashboard
// @Tags stats
// @Accept json
// @Produce json
// @Success 200 {object} string "dashboard stats"
// @Router /api/stats/dashboard [get]
func GetDashboardStats(c *gin.Context) {
	token, err := c.Cookie("admin_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization"})
		return
	}

	_, err = services.ValidateAdminAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	stats := services.GetDashboardStats()
	c.JSON(http.StatusOK, stats)
}

// @Summary Get Session Cache Stats
// @Description Returns information about the session cache
// @Tags stats
// @Accept json
// @Produce json
// @Success 200 {object} string "cache stats"
// @Router /api/stats/cache [get]
func GetCacheStats(c *gin.Context) {
	token, err := c.Cookie("admin_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorisation"})
		return
	}

	_, err = services.ValidateAdminAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	stats := services.GetCacheStats()
	c.JSON(http.StatusOK, stats)
}

// @Summary Get next session purge time
// @Description Get the specific time that the next invalid session purge will occur
// @Tags stats
// @Accept json
// @Produce json
// @Success 200 {object} string "time"
// @Router /api/stats/purge [get]
func GetNextPurgeTime(c *gin.Context) {
	token, err := c.Cookie("admin_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorisation"})
		return
	}

	_, err = services.ValidateAdminAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	purge_time := services.NextPurgeTime()
	c.JSON(http.StatusOK, purge_time)
}
