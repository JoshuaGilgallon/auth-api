package errors

// General errors
func NewNotFoundError(message string, err error) *UserError {
	return &UserError{Type: NotFound, Message: message, Err: err}
}

func NewValidationError(message string, err error) *UserError {
	return &UserError{Type: ValidationError, Message: message, Err: err}
}

func NewAlreadyExistsError(message string, err error) *UserError {
	return &UserError{Type: AlreadyExists, Message: message, Err: err}
}

func NewAuthenticationError(message string, err error) *UserError {
	return &UserError{Type: AuthenticationErr, Message: message, Err: err}
}

func NewInternalError(message string, err error) *UserError {
	return &UserError{Type: InternalError, Message: message, Err: err}
}

func NewFailedToEncryptError(message string, err error) *UserError {
	return &UserError{Type: FailedToEncrypt, Message: message, Err: err}
}

// Session errors
func NewSessionNotFoundError(message string, err error) *UserError {
	return &UserError{Type: SessionNotFound, Message: message, Err: err}
}

func NewSessionExpiredError(message string, err error) *UserError {
	return &UserError{Type: SessionExpired, Message: message, Err: err}
}

func NewInvalidTokenError(message string, err error) *UserError {
	return &UserError{Type: InvalidToken, Message: message, Err: err}
}

func NewTokenExpiredError(message string, err error) *UserError {
	return &UserError{Type: TokenExpired, Message: message, Err: err}
}

func NewMaxSessionsReachedError(message string, err error) *UserError {
	return &UserError{Type: MaxSessionsReached, Message: message, Err: err}
}

func NewRateLimitExceededError(message string, err error) *UserError {
	return &UserError{Type: RateLimitExceeded, Message: message, Err: err}
}

func NewFailedToCreateError(message string, err error) *UserError {
	return &UserError{Type: FailedToCreate, Message: message, Err: err}
}

// Login errors
func NewInvalidCredentialsError(message string, err error) *LoginError {
	return &LoginError{Type: InvalidCredentials, Message: message, Err: err}
}

func NewAccountLockedError(message string, err error) *LoginError {
	return &LoginError{Type: AccountLocked, Message: message, Err: err}
}

func NewAccountDisabledError(message string, err error) *LoginError {
	return &LoginError{Type: AccountDisabled, Message: message, Err: err}
}

func NewTooManyAttemptsError(message string, err error) *LoginError {
	return &LoginError{Type: TooManyAttempts, Message: message, Err: err}
}
