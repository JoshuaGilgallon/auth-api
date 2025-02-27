package services

import (
	"auth-api/internal/repositories"
	"time"
)

func DeleteInvalidSessions() error {
	err := repositories.DeleteInvalidSessions()
	if err != nil {
		return err
	}
	return nil
}

var purgeTime time.Time

func UpdatePurgeTime(newTime time.Time) error {
	purgeTime = newTime
	return nil
}

func NextPurgeTime() time.Time {
	return purgeTime
}
