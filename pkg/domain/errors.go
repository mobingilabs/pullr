package domain

type ErrorKind string

// Error is the generic error type used across Pullr
type Error struct {
	Kind    ErrorKind `json:"kind"`
	Message string    `json:"message"`
	Details string    `json:"-"`
}

// Error returns the error message
func (err *Error) Error() string {
	return err.Message
}

// WithDetails sets error's details
func (err *Error) WithDetails(details string) *Error {
	err.Details = details
	return err
}

// Error kinds
const (
	ErrKindNotFound     ErrorKind = "ERR_NOT_FOUND"
	ErrKindUnexpected   ErrorKind = "ERR_UNEXPECTED"
	ErrKindConflict     ErrorKind = "ERR_CONFLICT"
	ErrKindUnauthorized ErrorKind = "ERR_UNAUTHORIZED"
	ErrKindBadRequest   ErrorKind = "ERR_BADREQUEST"
	ErrKindUnsupported  ErrorKind = "ERR_UNSUPPORTED"
	ErrKindIrrelevant   ErrorKind = "ERR_IRRELEVANT"
)

var (
	// ErrNotFound is generic not found error
	ErrNotFound = &Error{ErrKindNotFound, "not found", ""}
	// ErrStorageDriver is generic storage driver error
	ErrStorageDriver = &Error{ErrKindUnexpected, "storage driver failed", ""}
	// ErrImageExists is image conflict error
	ErrImageExists = &Error{ErrKindConflict, "image exists", ""}
)

// AuthService errors
var (
	ErrAuthBadCredentials = &Error{ErrKindUnauthorized, "bad credentials", ""}
	ErrAuthUnauthorized   = &Error{ErrKindUnauthorized, "unauthenticated", ""}
	ErrAuthBadToken       = &Error{ErrKindUnauthorized, "invalid token", ""}
	ErrAuthTokenExpired   = &Error{ErrKindUnauthorized, "token expired", ""}
)

// OAuthService errors
var (
	ErrOAuthBadToken            = &Error{ErrKindBadRequest, "oauth: bad token", ""}
	ErrOAuthBadPayload          = &Error{ErrKindBadRequest, "oauth: bad payload", ""}
	ErrOAuthUnsupportedProvider = &Error{ErrKindUnsupported, "oauth: unsupported provider", ""}
)

// UserService errors
var (
	ErrUserUsernameExist = &Error{ErrKindConflict, "username is taken", ""}
	ErrUserEmailExist    = &Error{ErrKindConflict, "email is taken", ""}
)

// SourceService errors
var (
	ErrSourceUnsupportedProvider = &Error{ErrKindUnsupported, "unsupported source client", ""}
	ErrSourceBadPayload          = &Error{ErrKindBadRequest, "bad webhook payload", ""}
	ErrSourceIrrelevantEvent     = &Error{ErrKindIrrelevant, "irrelevant webhook event", ""}
)

// BuildService errors
var (
	ErrBuildBadJob = &Error{ErrKindBadRequest, "bad job", ""}
)
