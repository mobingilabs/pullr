package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// ErrorMiddleware turns Pullr errors to corresponding http errors
func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			pullrErr, ok := err.(*domain.Error)
			if !ok {
				return err
			}

			switch pullrErr {
			case domain.ErrNotFound:
				return echo.NewHTTPError(http.StatusNotFound, "not found")
			case domain.ErrStorageDriver:
				return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
			case domain.ErrAuthBadCredentials, domain.ErrAuthUnauthorized, domain.ErrAuthBadToken, domain.ErrAuthTokenExpired:
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			case domain.ErrOAuthBadToken:
				return echo.NewHTTPError(http.StatusBadRequest, "bad token")
			case domain.ErrOAuthBadPayload:
				return echo.NewHTTPError(http.StatusBadRequest, "bad payload")
			case domain.ErrOAuthUnsupportedProvider:
				return echo.NewHTTPError(http.StatusNotFound, "oauth provider not supported")
			case domain.ErrUserUsernameExist:
				return echo.NewHTTPError(http.StatusConflict, "username is taken")
			case domain.ErrUserEmailExist:
				return echo.NewHTTPError(http.StatusConflict, "email address is taken")
			}

			return err
		}
	}
}
