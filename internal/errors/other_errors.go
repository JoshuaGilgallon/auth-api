package errors

const (
	FailedToEncrypt ErrorType = "FAILED_TO_ENCRYPT"
)

func NewFailedToEncryptError(message string, err error) *UserError {
	return &UserError{
		Type:    FailedToEncrypt,
		Message: message,
		Err:     err,
	}
}
