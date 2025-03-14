package services

import (
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"auth-api/internal/utils"
	"fmt"
	"log"
	"strings"
	"time"
)

type UserInput struct {
	FirstName   string `json:"name"` // This field maps to "name" in JSON
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	BirthDate   string `json:"birth_date"`
	Language    string `json:"language"`
	MFAEnabled  bool   `json:"mfa_enabled"`
}

func CreateUser(input UserInput) (models.User, error) {
	// hash password before saving
	hashedPassword, err := utils.HashBcrypt(input.Password)

	if err != nil {
		return models.User{}, err
	}

	// encrypt the email and phone number before saving
	encryptedPhoneNumber, err := utils.Encrypt(input.PhoneNumber)
	if err != nil {
		return models.User{}, err
	}

	hashedPhoneNumber := utils.HashSHA(input.PhoneNumber)

	encryptedEmail, err := utils.Encrypt(input.Email)
	if err != nil {
		return models.User{}, err
	}

	hashedEmail := utils.HashSHA(input.Email)

	now := time.Now()
	user := models.User{
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		Email:           encryptedEmail,
		EmailHash:       hashedEmail,
		PhoneNumber:     encryptedPhoneNumber,
		PhoneNumberHash: hashedPhoneNumber,
		Password:        hashedPassword,
		Bio:             "",
		BirthDate:       time.Time{}, // Initialize with zero value
		Language:        input.Language,
		CreatedAt:       now,
		UpdatedAt:       now,
		MFAEnabled:      input.MFAEnabled,
		Status:          models.StatusActive, // set the status to active by default
	}

	if input.BirthDate != "" {
		birthDate, err := time.Parse("2006-01-02", input.BirthDate)
		if err != nil {
			return models.User{}, err
		}
		user.BirthDate = birthDate
	}

	return repositories.SaveUser(user)
}

func UpdateUser(id string, input UserInput) (models.User, error) {
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
	if input.Email != "" {
		// Encrypt and hash email
		encryptedEmail, err := utils.Encrypt(input.Email)
		if err != nil {
			log.Printf("Error encrypting email for user %s: %v", id, err)
			return models.User{}, fmt.Errorf("failed to encrypt email: %w", err)
		}
		existingUser.Email = encryptedEmail
		existingUser.EmailHash = utils.HashSHA(strings.ToLower(input.Email)) // Hash lowercase email
	}
	if input.PhoneNumber != "" {
		// Encrypt and hash phone number
		encryptedPhone, err := utils.Encrypt(input.PhoneNumber)
		if err != nil {
			log.Printf("Error encrypting phone for user %s: %v", id, err)
			return models.User{}, fmt.Errorf("failed to encrypt phone: %w", err)
		}
		existingUser.PhoneNumber = encryptedPhone
		existingUser.PhoneNumberHash = utils.HashSHA(input.PhoneNumber)
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

	// decrypt email
	decryptedEmail, err := utils.Decrypt(user.Email)
	if err != nil {
		return models.User{}, err
	}
	user.Email = decryptedEmail

	// decrypt phone number
	decryptedPhone, err := utils.Decrypt(user.PhoneNumber)
	if err != nil {
		return models.User{}, err
	}
	user.PhoneNumber = decryptedPhone

	return user, nil
}

func SearchUserByCredentials(search_term string) (models.User, error) {
	// get user with matching email
	user, err := repositories.GetUserByEmail(search_term)
	if err != nil {
		// if user not found with email, try searching by phone number
		user, err = repositories.GetUserByPhoneNumber(search_term)
		if err != nil {
			return models.User{}, err
		}
	}

	// decrypt email
	decryptedEmail, err := utils.Decrypt(user.Email)
	if err != nil {
		return models.User{}, err
	}

	user.Email = decryptedEmail

	// decrypt phone number
	decryptedPhoneNumber, err := utils.Decrypt(user.PhoneNumber)
	if err != nil {
		return models.User{}, err
	}

	user.PhoneNumber = decryptedPhoneNumber

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

	var users []models.User
	for _, user_raw := range user_raw {
		decrypted_email, err := utils.Decrypt(user_raw.Email)
		if err != nil {
			return SearchResult{}, err
		}

		decrypted_phone_number, err := utils.Decrypt(user_raw.PhoneNumber)
		if err != nil {
			return SearchResult{}, err
		}

		new_user_model := models.User{
			ID:              user_raw.ID,
			FirstName:       user_raw.FirstName,
			LastName:        user_raw.LastName,
			Email:           decrypted_email,
			EmailHash:       user_raw.EmailHash,
			PhoneNumber:     decrypted_phone_number,
			PhoneNumberHash: user_raw.PhoneNumberHash,
			Password:        user_raw.Password,
			Bio:             user_raw.Bio,
			BirthDate:       user_raw.BirthDate,
			Language:        user_raw.Language,
			CreatedAt:       user_raw.CreatedAt,
			UpdatedAt:       user_raw.UpdatedAt,
			MFAEnabled:      user_raw.MFAEnabled,
			Status:          user_raw.Status,
		}
		users = append(users, new_user_model)
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

	var users []models.User
	for _, user_raw := range user_raw {
		decrypted_email, err := utils.Decrypt(user_raw.Email)
		if err != nil {
			return SearchResult{}, err
		}

		decrypted_phone_number, err := utils.Decrypt(user_raw.PhoneNumber)
		if err != nil {
			return SearchResult{}, err
		}

		new_user_model := models.User{
			ID:              user_raw.ID,
			FirstName:       user_raw.FirstName,
			LastName:        user_raw.LastName,
			Email:           decrypted_email,
			EmailHash:       user_raw.EmailHash,
			PhoneNumber:     decrypted_phone_number,
			PhoneNumberHash: user_raw.PhoneNumberHash,
			Password:        user_raw.Password,
			Bio:             user_raw.Bio,
			BirthDate:       user_raw.BirthDate,
			Language:        user_raw.Language,
			CreatedAt:       user_raw.CreatedAt,
			UpdatedAt:       user_raw.UpdatedAt,
			MFAEnabled:      user_raw.MFAEnabled,
			Status:          user_raw.Status,
		}
		users = append(users, new_user_model)
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

	// decrypt email
	decryptedEmail, err := utils.Decrypt(user.Email)
	if err != nil {
		return models.User{}, err
	}
	user.Email = decryptedEmail

	// decrypt phone number
	decryptedPhoneNumber, err := utils.Decrypt(user.PhoneNumber)
	if err != nil {
		return models.User{}, err
	}
	user.PhoneNumber = decryptedPhoneNumber

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
	if criteria.Email != "" || criteria.PhoneNumber != "" || criteria.FirstName != "" || criteria.LastName != "" {
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
			// Decrypt sensitive fields
			decryptedEmail, err := utils.Decrypt(user.Email)
			if err != nil {
				return SearchResult{}, err
			}

			decryptedPhone, err := utils.Decrypt(user.PhoneNumber)
			if err != nil {
				return SearchResult{}, err
			}

			user.Email = decryptedEmail
			user.PhoneNumber = decryptedPhone

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
	if criteria.PhoneNumber != "" && !strings.Contains(user.PhoneNumber, criteria.PhoneNumber) {
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

	// Decrypt sensitive information for each user
	for i := range users {
		// Decrypt email
		decryptedEmail, err := utils.Decrypt(users[i].Email)
		if err != nil {
			return SearchResult{}, err
		}
		users[i].Email = decryptedEmail

		// Decrypt phone number
		decryptedPhone, err := utils.Decrypt(users[i].PhoneNumber)
		if err != nil {
			return SearchResult{}, err
		}
		users[i].PhoneNumber = decryptedPhone
	}

	return SearchResult{
		Users:        users,
		TotalResults: total,
	}, nil
}
