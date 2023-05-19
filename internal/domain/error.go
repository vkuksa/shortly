package shortly

import (
	"errors"
	"fmt"
)

const (
	ErrInternal = "internal"
	ErrNotFound = "not_found"
	ErrInvalid  = "invalid"
	ErrConflict = "already_exists"
)

// Error represents an application-specific error.
type Error struct {
	// Machine-readable error code.
	Code string

	// Human-readable error message.
	Message string
}

// Error implements the error interface. Not used by the application otherwise.
func (e *Error) Error() string {
	return fmt.Sprintf("shortly error: code=%s message=%s", e.Code, e.Message)
}

// Errorf is a helper function to return an Error with a given code and formatted message.
func NewError(code string, format string, args ...interface{}) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// ErrorCode unwraps an application error and returns its code.
// Non-application errors always return EINTERNAL.
func ErrorCode(err error) string {
	var e *Error

	switch {
	case err == nil:
		return ""
	case errors.As(err, &e):
		return e.Code
	default:
		return ErrInternal
	}
}

// ErrorMessage unwraps an application error and returns its message.
// Non-application errors always return "Internal error".
func ErrorMessage(err error) string {
	var e *Error

	switch {
	case err == nil:
		return ""
	case errors.As(err, &e):
		return e.Message
	default:
		return "Internal error."
	}
}
