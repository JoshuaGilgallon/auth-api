package services

import (
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"time"

	"auth-api/internal/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
