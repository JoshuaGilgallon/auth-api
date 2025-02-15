package errors

const (
	SessionNotFound    ErrorType = "SESSION_NOT_FOUND"
	SessionExpired     ErrorType = "SESSION_EXPIRED"
	InvalidToken       ErrorType = "INVALID_TOKEN"
	TokenExpired       ErrorType = "TOKEN_EXPIRED"
	MaxSessionsReached ErrorType = "MAX_SESSIONS_REACHED"
	RateLimitExceeded  ErrorType = "RATE_LIMIT_EXCEEDED"
	FailedToCreate     ErrorType = "FAILED_TO_CREATE"
)

// Helper functions for session-specific errors
func NewSessionNotFoundError(message string, err error) *UserError {
	return &UserError{
		Type:    SessionNotFound,
		Message: message,
		Err:     err,
	}
}

func NewSessionExpiredError(message string, err error) *UserError {
	return &UserError{
		Type:    SessionExpired,
		Message: message,
		Err:     err,
	}
}

func NewInvalidTokenError(message string, err error) *UserError {
	return &UserError{
		Type:    InvalidToken,
		Message: message,
		Err:     err,
	}
}

func NewTokenExpiredError(message string, err error) *UserError {
	return &UserError{
		Type:    TokenExpired,
		Message: message,
		Err:     err,
	}
}

func NewMaxSessionsReachedError(message string, err error) *UserError {
	return &UserError{
		Type:    MaxSessionsReached,
		Message: message,
		Err:     err,
	}
}

func NewRateLimitExceededError(message string, err error) *UserError {
	return &UserError{
		Type:    RateLimitExceeded,
		Message: message,
		Err:     err,
	}
}

func NewFailedToCreateError(message string, err error) *UserError {
	return &UserError{
		Type:    FailedToCreate,
		Message: message,
		Err:     err,
	}
}
