package services

import (
	"auth-api/internal/errors"
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"auth-api/internal/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminLoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateAdminSession(adminID primitive.ObjectID) (models.AdminSession, error) {
	if !limiter.Allow() {
		return models.AdminSession{}, errors.NewRateLimitExceededError("rate limit exceeded", nil)
	}

	sessionLock.Lock()
	defer sessionLock.Unlock()

	activeSessions, err := repositories.GetActiveSessionsByUserID(adminID)
	if err != nil {
		return models.AdminSession{}, errors.NewInternalError("failed to check active sessions", err)
	}

	if len(activeSessions) >= maxSessionsPerUser {
		return models.AdminSession{}, errors.NewMaxSessionsReachedError("maximum concurrent sessions reached", nil)
	}

	accessToken, err := generateSecureToken()
	if err != nil {
		return models.AdminSession{}, err
	}

	now := time.Now()
	session := models.AdminSession{
		AdminID:         adminID,
		AccessToken:     accessToken,
		AccessExpiresAt: now.Add(accessTokenDuration),
		CreatedAt:       now,
		LastActivity:    now,
	}

	session, err = repositories.SaveAdminSession(session)
	if err != nil {
		return models.AdminSession{}, errors.NewInternalError("failed to save session", err)
	}

	return session, nil
}

func ValidateAdminAccessToken(accessToken string) (models.AdminSession, error) {
	adminSession, err := repositories.GetAdminSessionByAccessToken(accessToken)
	if err != nil {
		return models.AdminSession{}, errors.NewSessionNotFoundError("session not found", err)
	}

	now := time.Now()

	// Check if session has expired due to inactivity or max lifespan
	if now.Sub(adminSession.LastActivity) > refreshTokenDuration {
		repositories.DeleteAdminSession(adminSession.AdminID)
		return models.AdminSession{}, errors.NewSessionExpiredError("session expired due to inactivity", nil)
	}

	if now.Sub(adminSession.CreatedAt) > maxSessionLifespan {
		repositories.DeleteAdminSession(adminSession.AdminID)
		return models.AdminSession{}, errors.NewSessionExpiredError("session exceeded maximum lifespan", nil)
	}

	// Check if access token has expired
	if now.After(adminSession.AccessExpiresAt) {
		return models.AdminSession{}, errors.NewTokenExpiredError("access token expired", nil)
	}

	// Update last activity
	adminSession.LastActivity = now
	repositories.UpdateAdminSession(adminSession)

	return adminSession, nil
}

func AdminLogin(input AdminLoginInput) (models.AdminSession, error) {
	var adminUser models.AdminUser

	// Try to get the user by Email first
	adminUser, err := repositories.GetAdminByUsername(input.Username)
	if adminUser.ID.IsZero() {
		log.Printf("Admin account not found with %v", err)
	}

	// Validate password
	if !utils.ValidateBcrypt(input.Password, adminUser.Password) {
		return models.AdminSession{}, errors.NewInvalidCredentialsError("The Password is Incorrect", nil)
	}

	// Create session
	adminSession, err := CreateAdminSession(adminUser.ID)
	if err != nil {
		return models.AdminSession{}, errors.NewFailedToCreateError("Failed to create session. Please wait a moment and try again, or contact support for assistance.", nil)
	}

	return adminSession, nil
}

func AdminLogout(accessToken string) error {
	return repositories.InvalidateAdminSessionByAccessToken(accessToken)
}
