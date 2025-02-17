package services

import (
	"auth-api/internal/models"
	"auth-api/internal/repositories"
	"auth-api/internal/utils"
	"time"
)

type UserInput struct {
	FirstName   string `json:"name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
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
		CreatedAt:       now,
		UpdatedAt:       now,
		MFAEnabled:      input.MFAEnabled,
		Status:          models.StatusActive, // set the status to active by default
	}
	return repositories.SaveUser(user)
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

func SearchUserByCreateTimeRange(start, end time.Time) ([]models.User, error) {
	user_raw, err := repositories.GetUsersByTimeCreatedRange(start, end)
	if err != nil {
		return nil, err
	}

	var users []models.User
	for _, user_raw := range user_raw {
		decrypted_email, err := utils.Decrypt(user_raw.Email)
		if err != nil {
			return nil, err
		}

		decrypted_phone_number, err := utils.Decrypt(user_raw.PhoneNumber)
		if err != nil {
			return nil, err
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
			CreatedAt:       user_raw.CreatedAt,
			UpdatedAt:       user_raw.UpdatedAt,
			MFAEnabled:      user_raw.MFAEnabled,
			Status:          user_raw.Status,
		}
		users = append(users, new_user_model)
	}

	return users, nil
}

func SearchUsersByTimeUpdatedRange(start, end time.Time) ([]models.User, error) {
	user_raw, err := repositories.GetUsersByTimeUpdatedRange(start, end)
	if err != nil {
		return nil, err
	}

	var users []models.User
	for _, user_raw := range user_raw {
		decrypted_email, err := utils.Decrypt(user_raw.Email)
		if err != nil {
			return nil, err
		}

		decrypted_phone_number, err := utils.Decrypt(user_raw.PhoneNumber)
		if err != nil {
			return nil, err
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
			CreatedAt:       user_raw.CreatedAt,
			UpdatedAt:       user_raw.UpdatedAt,
			MFAEnabled:      user_raw.MFAEnabled,
			Status:          user_raw.Status,
		}
		users = append(users, new_user_model)
	}

	return users, nil
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
