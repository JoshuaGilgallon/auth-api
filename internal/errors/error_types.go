package errors

type ErrorType string

// Error implements error.
func (e ErrorType) Error() string {
	panic("unimplemented")
}

type LoginErrorType string

const (
	// General errors
	NotFound          ErrorType = "NOT_FOUND"
	ValidationError   ErrorType = "VALIDATION_ERROR"
	AlreadyExists     ErrorType = "ALREADY_EXISTS"
	AuthenticationErr ErrorType = "AUTHENTICATION_ERROR"
	InternalError     ErrorType = "INTERNAL_ERROR"
	FailedToEncrypt   ErrorType = "FAILED_TO_ENCRYPT"

	// Session errors
	SessionNotFound    ErrorType = "SESSION_NOT_FOUND"
	SessionExpired     ErrorType = "SESSION_EXPIRED"
	InvalidToken       ErrorType = "INVALID_TOKEN"
	TokenExpired       ErrorType = "TOKEN_EXPIRED"
	MaxSessionsReached ErrorType = "MAX_SESSIONS_REACHED"
	RateLimitExceeded  ErrorType = "RATE_LIMIT_EXCEEDED"
	FailedToCreate     ErrorType = "FAILED_TO_CREATE"

	// Login specific errors
	InvalidCredentials LoginErrorType = "INVALID_CREDENTIALS"
	AccountLocked      LoginErrorType = "ACCOUNT_LOCKED"
	AccountDisabled    LoginErrorType = "ACCOUNT_DISABLED"
	TooManyAttempts    LoginErrorType = "TOO_MANY_ATTEMPTS"
)
