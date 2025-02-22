package services

import (
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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
	cacheCleanupInterval = 30 * time.Minute
)

var (
	limiter     *rate.Limiter
	sessionLock sync.RWMutex

	// Session cache:
	sessionCache    = make(map[string]models.Session)     // access token -> session
	refreshTokenMap = make(map[string]string)             // refresh token -> access token
	sessionIDMap    = make(map[primitive.ObjectID]string) // session ID -> access token
	cacheMutex      = sync.RWMutex{}
)

func init() {
	// Rate limit: 10 session creations per minute
	limiter = rate.NewLimiter(rate.Every(time.Minute/10), 1)

	go startCacheCleanup()
}

func generateSecureToken() (string, error) {
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	hash := sha256.New()
	hash.Write(b)
	hash.Write([]byte(time.Now().String()))

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func addToCache(session models.Session) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	sessionCache[session.AccessToken] = session
	refreshTokenMap[session.RefreshToken] = session.AccessToken
	sessionIDMap[session.ID] = session.AccessToken
}

func removeFromCache(session models.Session) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	delete(refreshTokenMap, session.RefreshToken)
	delete(sessionIDMap, session.ID)
	delete(sessionCache, session.AccessToken)
}

func getSessionFromCacheByAccessToken(accessToken string) (models.Session, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	session, found := sessionCache[accessToken]
	return session, found
}

func getSessionFromCacheByRefreshToken(refreshToken string) (models.Session, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	accessToken, found := refreshTokenMap[refreshToken]
	if !found {
		return models.Session{}, false
	}

	session, found := sessionCache[accessToken]
	return session, found
}

func getSessionFromCacheByID(sessionID primitive.ObjectID) (models.Session, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	accessToken, found := sessionIDMap[sessionID]
	if !found {
		return models.Session{}, false
	}

	session, found := sessionCache[accessToken]
	return session, found
}

func startCacheCleanup() {
	ticker := time.NewTicker(cacheCleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		cleanupExpiredSessions()
	}
}

func cleanupExpiredSessions() {
	now := time.Now()
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	for accessToken, session := range sessionCache {
		// remove if access token expired, session exceeded max lifespan, or inactive
		if now.After(session.AccessExpiresAt) ||
			now.Sub(session.CreatedAt) > maxSessionLifespan ||
			now.Sub(session.LastActivity) > refreshTokenDuration {
			delete(refreshTokenMap, session.RefreshToken)
			delete(sessionIDMap, session.ID)
			delete(sessionCache, accessToken)
		}
	}
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

	addToCache(session)

	return session, nil
}

func ValidateAccessToken(accessToken string) (models.Session, error) {
	// Try to get from cache first
	session, found := getSessionFromCacheByAccessToken(accessToken)

	// If not in cache try the database
	if !found {
		var err error
		session, err = repositories.GetSessionByAccessToken(accessToken)
		if err != nil {
			return models.Session{}, errors.NewSessionNotFoundError("session not found", err)
		}

		// add to cache for future requests
		addToCache(session)
	}

	now := time.Now()

	// check if session has expired due to inactivity or max lifespan
	if now.Sub(session.LastActivity) > refreshTokenDuration {
		repositories.DeleteSession(session.ID) // Remove from DB
		removeFromCache(session)               // Remove from cache
		return models.Session{}, errors.NewSessionExpiredError("session expired due to inactivity", nil)
	}

	if now.Sub(session.CreatedAt) > maxSessionLifespan {
		repositories.DeleteSession(session.ID) // Remove from DB
		removeFromCache(session)               // Remove from cache
		return models.Session{}, errors.NewSessionExpiredError("session exceeded maximum lifespan", nil)
	}

	// check if access token has expired
	if now.After(session.AccessExpiresAt) {
		return models.Session{}, errors.NewTokenExpiredError("access token expired", nil)
	}

	// update last activity
	session.LastActivity = now
	repositories.UpdateSession(session)

	addToCache(session)

	return session, nil
}

func RefreshAccessToken(refreshToken string) (models.Session, error) {
	// Try to get from cache first
	oldSession, found := getSessionFromCacheByRefreshToken(refreshToken)

	// If not in cache, try database
	if !found {
		var err error
		oldSession, err = repositories.GetSessionByRefreshToken(refreshToken)
		if err != nil {
			return models.Session{}, errors.NewSessionNotFoundError("session not found", err)
		}
	}

	now := time.Now()

	// Validate refresh token and session lifetime
	if now.After(oldSession.RefreshExpiresAt) ||
		now.Sub(oldSession.CreatedAt) > maxSessionLifespan ||
		now.Sub(oldSession.LastActivity) > refreshTokenDuration {
		repositories.DeleteSession(oldSession.ID) // Remove from DB
		removeFromCache(oldSession)               // Remove from cache
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

	// Remove old session from cache
	removeFromCache(oldSession)

	// Update session with new tokens
	oldSession.AccessToken = newAccessToken
	oldSession.RefreshToken = newRefreshToken
	oldSession.AccessExpiresAt = now.Add(accessTokenDuration)
	oldSession.RefreshExpiresAt = now.Add(refreshTokenDuration)
	oldSession.LastActivity = now

	updatedSession, err := repositories.UpdateSession(oldSession)
	if err != nil {
		return models.Session{}, err
	}

	// add the updated session to cache
	addToCache(updatedSession)

	return updatedSession, nil
}

func InvalidateSession(sessionID primitive.ObjectID) error {
	// Get it from the cache if available
	session, found := getSessionFromCacheByID(sessionID)
	if found {
		removeFromCache(session)
	}

	return repositories.DeleteSession(sessionID)
}

func InvalidateSessionByToken(token string) error {
	if token == "" {
		return errors.NewInvalidTokenError("empty token provided", nil)
	}

	// Try to get session by access token from cache first
	session, found := getSessionFromCacheByAccessToken(token)
	if found {
		removeFromCache(session)
		return repositories.DeleteSession(session.ID)
	}

	// Try to get session by refresh token from cache
	session, found = getSessionFromCacheByRefreshToken(token)
	if found {
		removeFromCache(session)
		return repositories.DeleteSession(session.ID)
	}

	// If not in cache, fall back to database operations:
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

func GetCacheStats() map[string]interface{} {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	return map[string]interface{}{
		"active_sessions": len(sessionCache),
	}
}
