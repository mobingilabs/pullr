package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/gova"
)

// ErrorMiddleware turns Pullr errors to corresponding http errors
func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			if pullrErr, ok := err.(*domain.Error); ok {
				return handlePullrError(pullrErr)
			}

			if validationErrs, ok := err.(gova.ValidationErrors); ok {
				return handleValidationErrors(validationErrs)
			}

			return err
		}
	}
}

func handlePullrError(err *domain.Error) error {
	switch err {
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

func handleValidationErrors(errs gova.ValidationErrors) error {
	response := make(map[string]string, len(errs))
	for _, err := range errs {
		response[err.Field] = err.Message
	}

	return echo.NewHTTPError(http.StatusBadRequest, response)
}
