package domain

// Error is error type used across Pullr
type Error struct {
	msg string
}

// Error returns error message
func (err *Error) Error() string {
	return err.msg
}

// NewError creates a pullr error
func NewError(msg string) *Error {
	return &Error{msg}
}

// ErrNotFound is generic not found error
var ErrNotFound = NewError("not found")

// ErrStorageDriver is generic storage driver error
var ErrStorageDriver = NewError("storage driver failed")

// AuthService errors
var (
	ErrAuthBadCredentials = NewError("bad credentials")
	ErrAuthUnauthorized   = NewError("unauthenticated")
	ErrAuthBadToken       = NewError("invalid token")
	ErrAuthTokenExpired   = NewError("token expired")
)

// OAuthService errors
var (
	ErrOAuthBadToken            = NewError("oauth: bad token")
	ErrOAuthBadPayload          = NewError("oauth: bad payload")
	ErrOAuthUnsupportedProvider = NewError("oauth: unsupported provider")
)

// UserService errors
var (
	ErrUserUsernameExist = NewError("username exist")
	ErrUserEmailExist    = NewError("email exist")
)
