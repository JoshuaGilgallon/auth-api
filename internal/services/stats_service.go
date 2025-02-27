package services

import (
	"auth-api/internal/repositories"
)

func GetCacheStats() map[string]interface{} {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	return map[string]interface{}{
		"active_sessions": len(sessionCache),
	}
}

func GetDashboardStats() map[string]interface{} {
	total_users, err := repositories.GetTotalUsers()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	total_sessions, err := repositories.GetTotalSessions()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	total_active_sessions, err := repositories.GetTotalActiveSessions()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	total_inactive_sessions, err := repositories.GetTotalInactiveSessions()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	weekly_logins, err := repositories.GetLoginCount()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	last_week_logins, err := repositories.GetLastWeekLoginCount()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	trends, err := repositories.GetWeeklyStats()
	if err != nil {
		return map[string]interface{}{
			"error": err.Error(),
		}
	}

	return map[string]interface{}{
		"total_users":             total_users,
		"total_sessions":          total_sessions,
		"total_active_sessions":   total_active_sessions,
		"total_inactive_sessions": total_inactive_sessions,
		"weekly_logins":           weekly_logins,
		"last_week_logins":        last_week_logins,
		"trends":                  trends,
	}
}
