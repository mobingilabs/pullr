package v1

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/auth"
	"github.com/mobingilabs/pullr/pkg/storage"
)

func errorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		status := http.StatusInternalServerError
		errKind := "ERR_UNEXPECTED"

		switch err {

		// Auth errors
		case auth.ErrInvalidToken, auth.ErrUnauthenticated, auth.ErrTokenExpired:
			errKind = "ERR_LOGIN"
			status = http.StatusUnauthorized
		case auth.ErrCredentials:
			errKind = "ERR_CREDENTIALS"
			status = http.StatusUnauthorized
		case auth.ErrUsernameTaken:
			errKind = "ERR_USERNAMETAKEN"
			status = http.StatusConflict

		// Storage errors
		case storage.ErrNotFound:
			errKind = "ERR_RESOURCE_NOTFOUND"
			status = http.StatusNotFound

		default:
			return err
		}

		return c.JSON(status, errMsg{errKind, ""})
	}
}
