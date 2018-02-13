package srv

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/auth"
	"github.com/mobingilabs/pullr/pkg/storage"
	"github.com/mobingilabs/pullr/pkg/vcs"
)

// ErrMsg describes a server error
type ErrMsg struct {
	Kind   string `json:"kind"`
	Status int    `json:"status"`
	Msg    string `json:"msg,omitempty"`
}

// NewErr creates an ErrMsg with given data
func NewErr(kind string, status int, msg string) ErrMsg {
	return ErrMsg{kind, status, msg}
}

// NewErrInternal creates an ErrMsg for internal server errors
func NewErrInternal() ErrMsg {
	return ErrMsg{"ERR_INTERNAL", http.StatusInternalServerError, "Unexpected error happened"}
}

// NewErrMissingParam creates an ErrMsg for a missing parameter
func NewErrMissingParam(param string) ErrMsg {
	msg := fmt.Sprintf("Query parameter '%s' is missing", param)
	return NewErr("ERR_MISSING_PARAM", http.StatusBadRequest, msg)
}

// NewErrBadValue creates an ErrMsg for a bad value
func NewErrBadValue(param, value string) ErrMsg {
	msg := fmt.Sprintf("Bad value '%s' for param '%s'", param, value)
	return NewErr("ERR_BAD_VALUE", http.StatusBadRequest, msg)
}

// NewErrUnsupported creates an ErrMsg for an unsupported feature
func NewErrUnsupported(feature string, vals ...interface{}) ErrMsg {
	var msg string
	if len(vals) > 0 {
		msg = fmt.Sprintf("%s is not supported", fmt.Sprintf(feature, vals))
	} else {
		msg = fmt.Sprintf("%s is not supported", feature)
	}

	return NewErr("ERR_UNSUPPORTED", http.StatusBadRequest, msg)
}

// Error reports error message
func (e ErrMsg) Error() string {
	return e.Msg
}

// ErrorHandler is an echo middleware to map few known error values from common
// packages as well as ErrMsg values.
func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		if e, ok := err.(ErrMsg); ok {
			return c.JSON(e.Status, e)
		}

		e := ErrMsg{}

		switch err {

		// Auth errors
		case auth.ErrInvalidToken, auth.ErrUnauthenticated, auth.ErrTokenExpired:
			e.Kind = "ERR_LOGIN"
			e.Status = http.StatusUnauthorized
			e.Msg = "Authentication required"
		case auth.ErrCredentials:
			e.Kind = "ERR_CREDENTIALS"
			e.Status = http.StatusUnauthorized
			e.Msg = "Wrong password or username"
		case auth.ErrUsernameTaken:
			e.Kind = "ERR_USERNAMETAKEN"
			e.Status = http.StatusConflict
			e.Msg = "Username is taken by another user"

		// Service errors
		case storage.ErrNotFound:
			e.Kind = "ERR_RESOURCE_NOTFOUND"
			e.Status = http.StatusNotFound
			e.Msg = "Resource not found"

		// Vcs errors
		case vcs.ErrInvalidWebhook, vcs.ErrInvalidWebhookPayload:
			e.Kind = "ERR_INVALID_WEBHOOK"
			e.Status = http.StatusBadRequest
			e.Msg = "Invalid webhook request"

		default:
			return err
		}

		return c.JSON(e.Status, e)
	}
}
