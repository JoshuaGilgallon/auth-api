package utils

import (
	"auth-api/internal/models"
	"regexp"
	"strings"
	"time"
)

func ExtractBearerToken(auth string) string {
	if auth == "" {
		return ""
	}
	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 32 {
		return false
	}
	match, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", username)
	return match
}

func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasNumber := strings.ContainsAny(password, "0123456789")
	hasSpecial := strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?")
	return hasUpper && hasLower && hasNumber && hasSpecial
}

func SafeParseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}

func DeduplicateUsers(users []models.User) []models.User {
	userMap := make(map[string]models.User)
	for _, user := range users {
		userMap[user.ID.Hex()] = user
	}

	uniqueUsers := make([]models.User, 0, len(userMap))
	for _, user := range userMap {
		uniqueUsers = append(uniqueUsers, user)
	}
	return uniqueUsers
}
