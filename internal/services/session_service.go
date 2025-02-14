package services

import (
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/time/rate"
)

const (
	sessionDuration      = 1 * time.Hour
	refreshTokenDuration = 24 * time.Hour
	tokenLength          = 32
	maxSessionsPerUser   = 5                    // Maximum concurrent sessions per user
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

func CreateSession(userID primitive.ObjectID, r *http.Request) (models.Session, error) {
	if !limiter.Allow() {
		return models.Session{}, errors.New("rate limit exceeded")
	}

	sessionLock.Lock()
	defer sessionLock.Unlock()

	// Check concurrent sessions
	activeSessions, _ := repositories.GetActiveSessionsByUserID(userID)
	if len(activeSessions) >= maxSessionsPerUser {
		return models.Session{}, errors.New("maximum sessions reached")
	}

	token, err := generateSecureToken()
	if err != nil {
		return models.Session{}, err
	}

	refreshToken, err := generateSecureToken()
	if err != nil {
		return models.Session{}, err
	}

	session := models.Session{
		UserID:       userID,
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(sessionDuration),
		IssuedAt:     time.Now(),
		IPAddress:    r.RemoteAddr,
		UserAgent:    r.UserAgent(),
		LastActivity: time.Now(),
	}

	return repositories.SaveSession(session)
}

func ValidateSession(token string, r *http.Request) (models.Session, error) {
	session, err := repositories.GetSessionByToken(token)
	if err != nil {
		return models.Session{}, err
	}

	// Validate session
	if time.Now().After(session.ExpiresAt) {
		repositories.DeleteSession(token)
		return models.Session{}, errors.New("session expired")
	}

	// Validate IP and User Agent for security
	if session.IPAddress != r.RemoteAddr {
		repositories.DeleteSession(token)
		return models.Session{}, errors.New("invalid session source")
	}

	// Update last activity
	session.LastActivity = time.Now()
	repositories.UpdateSession(session)

	return session, nil
}

func RefreshSession(refreshToken string, r *http.Request) (models.Session, error) {
	oldSession, err := repositories.GetSessionByRefreshToken(refreshToken)
	if err != nil {
		return models.Session{}, err
	}

	if time.Now().After(oldSession.ExpiresAt.Add(refreshTokenDuration)) {
		repositories.DeleteSession(oldSession.Token)
		return models.Session{}, errors.New("refresh token expired")
	}

	// Create new session
	newSession, err := CreateSession(oldSession.UserID, r)
	if err != nil {
		return models.Session{}, err
	}

	// Invalidate old session
	repositories.DeleteSession(oldSession.Token)

	return newSession, nil
}

func InvalidateSession(token string) error {
	return repositories.DeleteSession(token)
}
