package domain

import (
	"fmt"
)

// Error is the generic error type used across Pullr
type Error struct {
	msg string
}

// Error returns the error message
func (err *Error) Error() string {
	return err.msg
}

// ValidationError is data validation error type
type ValidationError struct {
	field string
	msg   string
}

// Error return the error message
func (err *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", err.field, err.msg)
}

// Validator is an helper for validating struct fields
type Validator struct {
	valid  bool
	errors []ValidationError
}

// Valid reports true if no errors encountered
func (v *Validator) Valid() bool {
	return v.valid
}

// Errors reports validation errors
func (v *Validator) Errors() []ValidationError {
	return v.errors
}

// NotEmptyString checks if the given string is not empty
func (v *Validator) NotEmptyString(field string, value string) {
	if value == "" {
		v.valid = false
		v.errors = append(v.errors, ValidationError{field, "can not be empty"})
	}
}

// NotEmpty checks if the given slice is not empty
func (v *Validator) NotEmpty(field string, nitems int) {
	if nitems == 0 {
		v.valid = false
		v.errors = append(v.errors, ValidationError{field, "should have at least one value"})
	}
}

// NonZero checks if the given integer is not zero
func (v *Validator) NonZero(field string, value int) {
	if value == 0 {
		v.valid = false
		v.errors = append(v.errors, ValidationError{field, "can not be zero"})
	}
}

// ErrNotFound is generic not found error
var ErrNotFound = &Error{"not found"}

// ErrStorageDriver is generic storage driver error
var ErrStorageDriver = &Error{"storage driver failed"}

// AuthService errors
var (
	ErrAuthBadCredentials = &Error{"bad credentials"}
	ErrAuthUnauthorized   = &Error{"unauthenticated"}
	ErrAuthBadToken       = &Error{"invalid token"}
	ErrAuthTokenExpired   = &Error{"token expired"}
)

// OAuthService errors
var (
	ErrOAuthBadToken            = &Error{"oauth: bad token"}
	ErrOAuthBadPayload          = &Error{"oauth: bad payload"}
	ErrOAuthUnsupportedProvider = &Error{"oauth: unsupported provider"}
)

// UserService errors
var (
	ErrUserUsernameExist = &Error{"username exist"}
	ErrUserEmailExist    = &Error{"email exist"}
)
