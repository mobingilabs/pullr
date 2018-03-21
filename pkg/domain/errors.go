package domain

// Error is the generic error type used across Pullr
type Error struct {
	msg string
}

// Error returns the error message
func (err *Error) Error() string {
	return err.msg
}

var (
	// ErrNotFound is generic not found error
	ErrNotFound = &Error{"not found"}
	// ErrStorageDriver is generic storage driver error
	ErrStorageDriver = &Error{"storage driver failed"}
	// ErrImageExists is image conflict error
	ErrImageExists = &Error{"image exists"}
)

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

// SourceService errors
var (
	ErrSourceUnsupportedProvider = &Error{"unsupported source client"}
	ErrSourceBadPayload          = &Error{"bad webhook payload"}
	ErrSourceIrrelevantEvent     = &Error{"irrelevant webhook event"}
)

// BuildService errors
var (
	ErrBuildBadJob = &Error{"bad job"}
)
