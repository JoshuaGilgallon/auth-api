package services

import (
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"auth-api/internal/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/time/rate"
)

const (
	accessTokenDuration  = 30 * time.Minute
	refreshTokenDuration = 12 * time.Hour
	maxSessionLifespan   = 7 * 24 * time.Hour // 7 days
	tokenLength          = 32
	maxSessionsPerUser   = 5 // Maximum concurrent sessions per user
)

var (
	limiter     *rate.Limiter
	sessionLock sync.RWMutex
)

func init() {
	// Rate limit: 10 session creations per minute
	limiter = rate.NewLimiter(rate.Every(time.Minute/10), 1)
}

func generateSecureToken() (string, error) {
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	hash := sha256.New()
	hash.Write(b)
	hash.Write([]byte(time.Now().String()))

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func CreateSession(userID primitive.ObjectID) (models.Session, error) {
	if !limiter.Allow() {
		return models.Session{}, errors.NewRateLimitExceededError("rate limit exceeded", nil)
	}

	sessionLock.Lock()
	defer sessionLock.Unlock()

	activeSessions, err := repositories.GetActiveSessionsByUserID(userID)
	if err != nil {
		return models.Session{}, errors.NewInternalError("failed to check active sessions", err)
	}

	if len(activeSessions) >= maxSessionsPerUser {
		return models.Session{}, errors.NewMaxSessionsReachedError("maximum concurrent sessions reached", nil)
	}

	accessToken, err := generateSecureToken()
	if err != nil {
		return models.Session{}, err
	}

	refreshToken, err := generateSecureToken()
	if err != nil {
		return models.Session{}, err
	}

	now := time.Now()
	session := models.Session{
		UserID:           userID,
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresAt:  now.Add(accessTokenDuration),
		RefreshExpiresAt: now.Add(refreshTokenDuration),
		CreatedAt:        now,
		LastActivity:     now,
	}

	session, err = repositories.SaveSession(session)
	if err != nil {
		return models.Session{}, errors.NewInternalError("failed to save session", err)
	}

	return session, nil
}

func ValidateAccessToken(accessToken string) (models.Session, error) {
	session, err := repositories.GetSessionByAccessToken(accessToken)
	if err != nil {
		return models.Session{}, errors.NewSessionNotFoundError("session not found", err)
	}

	now := time.Now()

	// Check if session has expired due to inactivity or max lifespan
	if now.Sub(session.LastActivity) > refreshTokenDuration {
		repositories.DeleteSession(session.ID)
		return models.Session{}, errors.NewSessionExpiredError("session expired due to inactivity", nil)
	}

	if now.Sub(session.CreatedAt) > maxSessionLifespan {
		repositories.DeleteSession(session.ID)
		return models.Session{}, errors.NewSessionExpiredError("session exceeded maximum lifespan", nil)
	}

	// Check if access token has expired
	if now.After(session.AccessExpiresAt) {
		return models.Session{}, errors.NewTokenExpiredError("access token expired", nil)
	}

	// Update last activity
	session.LastActivity = now
	repositories.UpdateSession(session)

	return session, nil
}

func RefreshAccessToken(refreshToken string) (models.Session, error) {
	oldSession, err := repositories.GetSessionByRefreshToken(refreshToken)
	if err != nil {
		return models.Session{}, errors.NewSessionNotFoundError("session not found", err)
	}

	now := time.Now()

	// Validate refresh token and session lifetime
	if now.After(oldSession.RefreshExpiresAt) ||
		now.Sub(oldSession.CreatedAt) > maxSessionLifespan ||
		now.Sub(oldSession.LastActivity) > refreshTokenDuration {
		repositories.DeleteSession(oldSession.ID)
		return models.Session{}, errors.NewTokenExpiredError("refresh token expired", nil)
	}

	// Generate new tokens
	newAccessToken, err := generateSecureToken()
	if err != nil {
		return models.Session{}, err
	}

	newRefreshToken, err := generateSecureToken()
	if err != nil {
		return models.Session{}, err
	}

	// Update session with new tokens
	oldSession.AccessToken = newAccessToken
	oldSession.RefreshToken = newRefreshToken
	oldSession.AccessExpiresAt = now.Add(accessTokenDuration)
	oldSession.RefreshExpiresAt = now.Add(refreshTokenDuration)
	oldSession.LastActivity = now

	return repositories.UpdateSession(oldSession)
}

func InvalidateSession(sessionID primitive.ObjectID) error {
	return repositories.DeleteSession(sessionID)
}

func InvalidateSessionByToken(token string) error {
	if token == "" {
		return errors.NewInvalidTokenError("empty token provided", nil)
	}

	// Try to get session by access token first
	session, err := repositories.GetSessionByAccessToken(token)
	if err == nil {
		return repositories.DeleteSession(session.ID)
	}

	// If not found by access token, try refresh token
	session, err = repositories.GetSessionByRefreshToken(token)
	if err == nil {
		return repositories.DeleteSession(session.ID)
	}

	return errors.NewSessionNotFoundError("invalid or expired session", nil)
}
