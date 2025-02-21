package services

import (
	"auth-api/internal/errors"
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"auth-api/internal/utils"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminLoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminCreationRequest struct {
	AdminUser AdminLoginInput `json:"user"`
	RootUser  AdminLoginInput `json:"root_user"`
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

func AdminLogin(input AdminLoginInput) (interface{}, error) {
	// Check if root user is logging in first
	if ValidateRootUserCredentials(input.Username, input.Password) {
		// Return a session-like structure for root user
		return map[string]interface{}{
			"access_token": "root",
			"is_root":      true,
		}, nil
	}

	// Try to get the user by username
	adminUser, err := repositories.GetAdminByUsername(input.Username)
	if err != nil || adminUser.ID.IsZero() {
		log.Printf("Admin account not found: %v", err)
		return nil, errors.NewInvalidCredentialsError("Invalid username or password", nil)
	}

	// Validate password
	if !utils.ValidateBcrypt(input.Password, adminUser.Password) {
		return nil, errors.NewInvalidCredentialsError("Incorrect password", nil)
	}

	// Create session for normal admin users
	adminSession, err := CreateAdminSession(adminUser.ID)
	if err != nil {
		return nil, errors.NewFailedToCreateError("Failed to create session", nil)
	}

	return map[string]interface{}{
		"access_token": adminSession.AccessToken,
		"is_root":      false,
	}, nil
}

func AdminLogout(accessToken string) error {
	return repositories.InvalidateAdminSessionByAccessToken(accessToken)
}

// root user stuff

func CreateAdminUser(input AdminLoginInput) (models.AdminUser, error) {
	// hash password before saving
	hashedPassword, err := utils.HashBcrypt(input.Password)
	if err != nil {
		return models.AdminUser{}, err
	}

	user := models.AdminUser{
		Username: input.Username,
		Password: hashedPassword,
	}
	return repositories.SaveAdmin(user)
}

func ValidateRootUserCredentials(username string, password string) bool {
	rootUsername := os.Getenv("ROOT_ADMIN_USERNAME")
	rootPassword := os.Getenv("ROOT_ADMIN_PASSWORD")

	// Debug prints
	log.Printf("Input username: '%s', env username: '%s'", username, rootUsername)
	log.Printf("Input password: '%s', env password: '%s'", password, rootPassword)

	if rootUsername == "" {
		log.Printf("Warning: ROOT_ADMIN_USERNAME environment variable is not set")
		return false
	}
	if rootPassword == "" {
		log.Printf("Warning: ROOT_ADMIN_PASSWORD environment variable is not set")
		return false
	}

	return username == rootUsername && password == rootPassword
}
