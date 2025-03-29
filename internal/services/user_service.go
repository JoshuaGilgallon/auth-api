package services

import (
	"auth-api/internal/errors"
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"auth-api/internal/utils"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(input models.UserInput) (models.User, error) {
	// hash password before saving
	hashedPassword, err := utils.HashBcrypt(input.Password)

	if err != nil {
		return models.User{}, err
	}

	// check if user already exists
	existingUser, err := repositories.GetUserByEmail(input.Email)
	if err == nil && existingUser.ID != primitive.NilObjectID {
		return models.User{}, errors.NewAlreadyExistsError("User with this email already exists", nil)
	}

	now := time.Now()

	user := models.User{
		FirstName:  "",
		LastName:   "",
		Email:      input.Email,
		Password:   hashedPassword,
		Bio:        "",
		BirthDate:  time.Time{}, // Initialize with zero value
		Language:   "EN_us",
		CreatedAt:  now,
		UpdatedAt:  now,
		MFAEnabled: false,
		Status:     models.StatusPending, // set the status to pending by default because they still have to verify and put all their other info in
	}

	// if input.BirthDate != "" {
	// 	birthDate, err := time.Parse("2006-01-02", input.BirthDate)
	// 	if err != nil {
	// 		return models.User{}, err
	// 	}
	// 	user.BirthDate = birthDate
	// }

	saved_user, err := repositories.SaveUser(user)
	if err != nil {
		log.Printf("Error saving user: %v", err)
		return models.User{}, fmt.Errorf("failed to save user: %w", err)
	}

	return saved_user, nil
}

func UpdateUser(id string, input models.FullUserInput) (models.User, error) {
	existingUser, err := repositories.GetUserByID(id)
	if err != nil {
		log.Printf("Error fetching user %s: %v", id, err)
		return models.User{}, fmt.Errorf("failed to fetch user: %w", err)
	}

	if input.FirstName != "" {
		existingUser.FirstName = input.FirstName
	}
	if input.LastName != "" {
		existingUser.LastName = input.LastName
	}
	if input.Password != "" {
		// Hash new password
		hashedPassword, err := utils.HashBcrypt(input.Password)
		if err != nil {
			log.Printf("Error hashing password for user %s: %v", id, err)
			return models.User{}, fmt.Errorf("failed to hash password: %w", err)
		}
		existingUser.Password = hashedPassword
	}
	if input.BirthDate != "" {
		birthDate, err := time.Parse("2006-01-02", input.BirthDate)
		if err != nil {
			log.Printf("Error parsing birth date '%s' for user %s: %v", input.BirthDate, id, err)
			return models.User{}, fmt.Errorf("invalid birth date format: %w", err)
		}
		existingUser.BirthDate = birthDate
	}
	if input.Language != "" {
		existingUser.Language = input.Language
	}

	existingUser.MFAEnabled = input.MFAEnabled
	existingUser.UpdatedAt = time.Now()

	// Save the updated user
	updatedUser, err := repositories.SaveUser(existingUser)
	if err != nil {
		log.Printf("Error saving user %s: %v", id, err)
		return models.User{}, fmt.Errorf("failed to save user: %w", err)
	}

	return updatedUser, nil
}

func GetUser(id string) (models.User, error) {
	user, err := repositories.GetUserByID(id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func SearchUserByCredentials(search_term string) (models.User, error) {
	// get user with matching email
	user, err := repositories.GetUserByEmail(search_term)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func SearchUserByCreateTimeRange(start, end time.Time, pageNumber, pageSize int64) (SearchResult, error) {
	if pageNumber <= 0 {
		pageNumber = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	skip := (pageNumber - 1) * pageSize
	user_raw, total, err := repositories.GetUsersByTimeCreatedRange(start, end, skip, pageSize)
	if err != nil {
		return SearchResult{}, err
	}

	users := make([]models.User, len(user_raw))
	for i, raw := range user_raw {
		users[i] = models.User(raw)
	}

	return SearchResult{
		Users:        users,
		TotalResults: total,
	}, nil
}

func SearchUsersByTimeUpdatedRange(start, end time.Time, pageNumber, pageSize int64) (SearchResult, error) {
	if pageNumber <= 0 {
		pageNumber = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	skip := (pageNumber - 1) * pageSize
	user_raw, total, err := repositories.GetUsersByTimeUpdatedRange(start, end, skip, pageSize)
	if err != nil {
		return SearchResult{}, err
	}

	users := make([]models.User, len(user_raw))
	for i, raw := range user_raw {
		users[i] = models.User(raw)
	}

	return SearchResult{
		Users:        users,
		TotalResults: total,
	}, nil
}
func GetCurrentUser(token string) (models.User, error) {
	session, err := repositories.GetSessionByAccessToken(token)
	if err != nil {
		return models.User{}, err
	}

	user, err := repositories.GetUserByID(session.UserID.Hex())
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func SearchUsers(criteria models.UserAdvancedSearchCriteria) (SearchResult, error) {
	// Set default pagination values if not provided
	if criteria.PageNumber <= 0 {
		criteria.PageNumber = 1
	}
	if criteria.PageSize <= 0 {
		criteria.PageSize = 10
	}

	skip := (criteria.PageNumber - 1) * criteria.PageSize
	var allUsers []models.User
	var totalResults int64 = 0

	// Search by time ranges if provided
	if criteria.StartTime != nil && criteria.EndTime != nil {
		users, total, err := repositories.GetUsersByTimeCreatedRange(
			*criteria.StartTime,
			*criteria.EndTime,
			skip,
			criteria.PageSize,
		)
		if err != nil {
			return SearchResult{}, err
		}
		allUsers = append(allUsers, users...)
		totalResults = total
	}

	if criteria.UpdateStartTime != nil && criteria.UpdateEndTime != nil {
		users, total, err := repositories.GetUsersByTimeUpdatedRange(
			*criteria.UpdateStartTime,
			*criteria.UpdateEndTime,
			skip,
			criteria.PageSize,
		)
		if err != nil {
			return SearchResult{}, err
		}
		allUsers = append(allUsers, users...)
		if total > totalResults {
			totalResults = total
		}
	}

	// Search by other criteria
	if criteria.Email != "" || criteria.FirstName != "" || criteria.LastName != "" {
		users, total, err := repositories.SearchUsersByFields(criteria)
		if err != nil {
			return SearchResult{}, err
		}
		allUsers = append(allUsers, users...)
		if total > totalResults {
			totalResults = total
		}
	}

	// Process results
	results := make([]models.User, 0)
	processedIDs := make(map[string]bool)

	for _, user := range allUsers {
		if !processedIDs[user.ID.Hex()] {
			// Add to results if matches criteria
			if matchesCriteria(user, criteria) {
				results = append(results, user)
				processedIDs[user.ID.Hex()] = true
			}
		}
	}

	return SearchResult{
		Users:        results,
		TotalResults: totalResults,
	}, nil
}

func matchesCriteria(user models.User, criteria models.UserAdvancedSearchCriteria) bool {
	if criteria.FirstName != "" && !strings.Contains(strings.ToLower(user.FirstName), strings.ToLower(criteria.FirstName)) {
		return false
	}
	if criteria.LastName != "" && !strings.Contains(strings.ToLower(user.LastName), strings.ToLower(criteria.LastName)) {
		return false
	}
	if criteria.Email != "" && !strings.Contains(strings.ToLower(user.Email), strings.ToLower(criteria.Email)) {
		return false
	}
	return true
}

type SearchResult struct {
	Users        []models.User
	TotalResults int64
}

func SimpleSearch(searchTerm string, pageNum int64, pageSize int64) (SearchResult, error) {
	skip := (pageNum - 1) * pageSize
	users, total, err := repositories.SimpleSearchUsers(searchTerm, skip, pageSize)
	if err != nil {
		return SearchResult{}, err
	}

	return SearchResult{
		Users:        users,
		TotalResults: total,
	}, nil
}
