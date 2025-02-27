package handlers

import (
	"auth-api/internal/errors"
	"auth-api/internal/services"
	"auth-api/internal/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Create session
// @Description Create a new session for a user
// @Tags sessions
// @Accept json
// @Produce json
// @Param user_id body string true "User ID"
// @Success 201 {object} models.Session
// @Router /api/session [post]
func CreateSession(c *gin.Context) {
	var input struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid session creation request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Sanitize input
	input.UserID = strings.TrimSpace(input.UserID)

	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		log.Printf("Invalid user ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	session, err := services.CreateSession(userID)
	if err != nil {
		log.Printf("Session creation failed for user %s: %v", input.UserID, err)

		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.RateLimitExceeded:
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many session requests"})
			case errors.MaxSessionsReached:
				c.JSON(http.StatusConflict, gin.H{"error": "Maximum active sessions reached"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
			}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	log.Printf("Session created successfully for user: %s", input.UserID)
	c.JSON(http.StatusCreated, session)
}

// @Summary Validate access token
// @Description Validate an access token and return session info
// @Tags sessions
// @Accept json
// @Produce json
// @Param authorization header string true "Bearer <token>"
// @Success 200 {object} models.Session
// @Router /api/session/validate [get]
func ValidateSession(c *gin.Context) {
	token := utils.ExtractBearerToken(c.GetHeader("Authorization"))
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	session, err := services.ValidateAccessToken(token)
	if err != nil {
		log.Printf("Session validation failed: %v", err)

		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.SessionExpired, errors.TokenExpired:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired"})
			case errors.SessionNotFound:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			case errors.InvalidToken:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token format"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Session validation failed"})
			}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, session)
}

// @Summary Refresh access token
// @Description Get new access and refresh tokens using a refresh token
// @Tags sessions
// @Accept json
// @Produce json
// @Param refresh_token body string true "Refresh Token"
// @Success 200 {object} models.Session
// @Router /api/session/refresh [post]
func RefreshSession(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid refresh token request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Sanitize input
	input.RefreshToken = strings.TrimSpace(input.RefreshToken)
	if input.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	session, err := services.RefreshAccessToken(input.RefreshToken)
	if err != nil {
		log.Printf("Session refresh failed: %v", err)

		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.TokenExpired:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
			case errors.SessionNotFound:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			case errors.InvalidToken:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token format"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh session"})
			}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	log.Printf("Session refreshed successfully")
	c.JSON(http.StatusOK, session)
}

// @Summary Invalidate session
// @Description Invalidate a session using its access token
// @Tags sessions
// @Accept json
// @Produce json
// @Param token path string true "Access Token"
// @Success 200 {object} string "session invalidated"
// @Router /api/session/{token} [delete]
func InvalidateSession(c *gin.Context) {
	token := utils.ExtractBearerToken(c.Param("session_id"))
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	if err := services.InvalidateSessionByToken(token); err != nil {
		log.Printf("Session invalidation failed: %v", err)

		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.SessionNotFound:
				c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
			case errors.InvalidToken:
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token format"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate session"})
			}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Session invalidated"})
}
