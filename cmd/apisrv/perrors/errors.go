package errors

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/auth"
	"github.com/mobingilabs/pullr/pkg/storage"
)

type errMsg struct {
	Kind   string `json:"kind"`
	Status int    `json:"status"`
	Msg    string `json:"msg,omitempty"`
}

func newErr(kind string, status int, msg string) errMsg {
	return errMsg{kind, status, msg}
}

func newErrMissingParam(param string) errMsg {
	msg := fmt.Sprintf("Query parameter '%s' is missing", param)
	return newErr("ERR_MISSING_PARAM", http.StatusBadRequest, msg)
}

func newErrBadValue(param, value string) errMsg {
	msg := fmt.Sprintf("Bad value '%s' for param '%s'", param, value)
	return newErr("ERR_BAD_VALUE", http.StatusBadRequest, msg)
}

func (e errMsg) Error() string {
	return e.Msg
}

func errorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		if e, ok := err.(errMsg); ok {
			return c.JSON(e.Status, e)
		}

		e := errMsg{}

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

		// Storage errors
		case storage.ErrNotFound:
			e.Kind = "ERR_RESOURCE_NOTFOUND"
			e.Status = http.StatusNotFound
			e.Msg = "Resource not found"

		// OAuth errors
		case ErrUnsupportedOAuthProvider:
			e.Kind = "ERR_OAUTH_UNSUPPORTEDPROVIDER"
			e.Status = http.StatusBadRequest
			e.Msg = "OAuth provider not supported"

		default:
			return err
		}

		return c.JSON(e.Status, e)
	}
}
