package errors

import "fmt"

type ErrorType string

const (
	NotFound          ErrorType = "NOT_FOUND"
	ValidationError   ErrorType = "VALIDATION_ERROR"
	AlreadyExists     ErrorType = "ALREADY_EXISTS"
	AuthenticationErr ErrorType = "AUTHENTICATION_ERROR"
	InternalError     ErrorType = "INTERNAL_ERROR"
)

type UserError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *UserError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *UserError) Unwrap() error {
	return e.Err
}

// Helper functions to create specific errors
func NewNotFoundError(message string, err error) *UserError {
	return &UserError{
		Type:    NotFound,
		Message: message,
		Err:     err,
	}
}

func NewValidationError(message string, err error) *UserError {
	return &UserError{
		Type:    ValidationError,
		Message: message,
		Err:     err,
	}
}

func NewAlreadyExistsError(message string, err error) *UserError {
	return &UserError{
		Type:    AlreadyExists,
		Message: message,
		Err:     err,
	}
}

func NewAuthenticationError(message string, err error) *UserError {
	return &UserError{
		Type:    AuthenticationErr,
		Message: message,
		Err:     err,
	}
}

func NewInternalError(message string, err error) *UserError {
	return &UserError{
		Type:    InternalError,
		Message: message,
		Err:     err,
	}
}