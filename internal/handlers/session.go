package handlers

import (
	"auth-api/internal/errors"
	"auth-api/internal/services"
	"net/http"

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
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid request format", err))
		return
	}

	userID, err := primitive.ObjectIDFromHex(input.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid user ID format", err))
		return
	}

	session, err := services.CreateSession(userID)
	if err != nil {
		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.RateLimitExceeded:
				c.JSON(http.StatusTooManyRequests, e)
			case errors.MaxSessionsReached:
				c.JSON(http.StatusConflict, e)
			default:
				c.JSON(http.StatusInternalServerError, e)
			}
		default:
			c.JSON(http.StatusInternalServerError, errors.NewInternalError("failed to create session", err))
		}
		return
	}

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
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.GetHeader("Token")
	}

	if token == "" {
		c.JSON(http.StatusBadRequest, errors.NewValidationError("no token provided", nil))
		return
	}

	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	session, err := services.ValidateAccessToken(token)
	if err != nil {
		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.SessionExpired, errors.TokenExpired:
				c.JSON(http.StatusUnauthorized, e)
			case errors.SessionNotFound:
				c.JSON(http.StatusNotFound, e)
			default:
				c.JSON(http.StatusInternalServerError, e)
			}
		default:
			c.JSON(http.StatusInternalServerError, errors.NewInternalError("failed to validate token", err))
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
		c.JSON(http.StatusBadRequest, errors.NewValidationError("invalid refresh token format", err))
		return
	}

	session, err := services.RefreshAccessToken(input.RefreshToken)
	if err != nil {
		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.TokenExpired:
				c.JSON(http.StatusUnauthorized, e)
			case errors.SessionNotFound:
				c.JSON(http.StatusNotFound, e)
			default:
				c.JSON(http.StatusInternalServerError, e)
			}
		default:
			c.JSON(http.StatusInternalServerError, errors.NewInternalError("failed to refresh token", err))
		}
		return
	}

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
	token := c.Param("session_id")
	if err := services.InvalidateSessionByToken(token); err != nil {
		switch e := err.(type) {
		case *errors.UserError:
			switch e.Type {
			case errors.SessionNotFound:
				c.JSON(http.StatusNotFound, e)
			case errors.InvalidToken:
				c.JSON(http.StatusBadRequest, e)
			default:
				c.JSON(http.StatusInternalServerError, e)
			}
		default:
			c.JSON(http.StatusInternalServerError, errors.NewInternalError("failed to invalidate session", err))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session invalidated"})
}
