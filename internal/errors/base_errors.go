package errors

import "log"

type UserError struct {
	Type    ErrorType
	Message string
	Err     error
}

type LoginError struct {
	Type    LoginErrorType
	Message string
	Err     error
}

func (e *UserError) Error() string {
	errorMessage := formatError(string(e.Type), e.Message, e.Err)
	log.Println(errorMessage)
	return errorMessage
}

func (e *LoginError) Error() string {
	errorMessage := formatError(string(e.Type), e.Message, e.Err)
	log.Println(errorMessage)
	return errorMessage
}

func (e *UserError) Unwrap() error {
	return e.Err
}

func (e *LoginError) Unwrap() error {
	return e.Err
}

func formatError(errType string, message string, err error) string {
	if err != nil {
		return errType + ": " + message + " (" + err.Error() + ")"
	}
	return errType + ": " + message
}
