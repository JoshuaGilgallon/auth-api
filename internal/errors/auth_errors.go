package errors

import "fmt"

type LoginErrorType string

const (
	InvalidCredentials LoginErrorType = "INVALID_CREDENTIALS"
	AccountLocked     LoginErrorType = "ACCOUNT_LOCKED"
	AccountDisabled   LoginErrorType = "Account_DISABLED"
	TooManyAttempts   LoginErrorType = "TOO_MANY_ATTEMPTS"
)

type LoginError struct {
	Type    LoginErrorType
	Message string
	Err     error
}

func (e *LoginError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *LoginError) Unwrap() error {
	return e.Err
}

// Helper functions to create specific login errors
func NewInvalidCredentialsError(message string, err error) *LoginError {
	return &LoginError{
		Type:    InvalidCredentials,
		Message: message,
		Err:     err,
	}
}

func NewAccountLockedError(message string, err error) *LoginError {
	return &LoginError{
		Type:    AccountLocked,
		Message: message,
		Err:     err,
	}
}

func NewAccountDisabledError(message string, err error) *LoginError {
	return &LoginError{
		Type:    AccountDisabled,
		Message: message,
		Err:     err,
	}
}

func NewTooManyAttemptsError(message string, err error) *LoginError {
	return &LoginError{
		Type:    TooManyAttempts,
		Message: message,
		Err:     err,
	}
}
