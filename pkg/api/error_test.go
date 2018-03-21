package api

import (
	"net/http"
	"testing"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

func TestErrorMiddleware(t *testing.T) {
	tests := map[error]int{
		domain.ErrNotFound:                 http.StatusNotFound,
		domain.ErrStorageDriver:            http.StatusInternalServerError,
		domain.ErrImageExists:              http.StatusConflict,
		domain.ErrAuthBadCredentials:       http.StatusUnauthorized,
		domain.ErrAuthUnauthorized:         http.StatusUnauthorized,
		domain.ErrAuthBadToken:             http.StatusUnauthorized,
		domain.ErrAuthTokenExpired:         http.StatusUnauthorized,
		domain.ErrOAuthBadToken:            http.StatusBadRequest,
		domain.ErrOAuthBadPayload:          http.StatusBadRequest,
		domain.ErrOAuthUnsupportedProvider: http.StatusNotFound,
		domain.ErrUserUsernameExist:        http.StatusConflict,
		domain.ErrUserEmailExist:           http.StatusConflict,
	}

	middleware := ErrorMiddleware()
	for err, status := range tests {
		newErr := middleware(func(c echo.Context) error { return err })(nil)

		httpErr, ok := newErr.(*echo.HTTPError)
		if !ok {
			t.Errorf(`"%v" error should be catched`, err)
		}

		if httpErr.Code != status {
			t.Errorf(`"%v" error should have status code %d`, err, status)
		}
	}
}
