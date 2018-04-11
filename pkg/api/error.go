package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/gova"
)

// ErrorMiddleware turns Pullr errors to corresponding http errors
func ErrorMiddleware(logger domain.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			if pullrErr, ok := err.(*domain.Error); ok {
				return handlePullrError(logger, c, pullrErr)
			}
			if validationErrs, ok := err.(gova.ValidationErrors); ok {
				return handleValidationErrors(validationErrs)
			}

			logger.Errorf("unexpected error: %v", err)
			return err
		}
	}
}

func handlePullrError(logger domain.Logger, c echo.Context, err *domain.Error) error {
	switch err.Kind {
	case domain.ErrKindNotFound:
		return echo.ErrNotFound
		return c.JSON(http.StatusNotFound, err)
	case domain.ErrKindUnexpected:
		logger.Errorf("%v: %s", err, err.Details)
		return c.JSON(http.StatusInternalServerError, err)
	case domain.ErrKindConflict:
		return c.JSON(http.StatusConflict, err)
	case domain.ErrKindUnauthorized:
		return c.JSON(http.StatusUnauthorized, err)
	case domain.ErrKindBadRequest:
		return c.JSON(http.StatusBadRequest, err)
	}

	logger.Errorf("pullr: %v", err)
	return err
}

func handleValidationErrors(errs gova.ValidationErrors) error {
	response := make(map[string]string, len(errs))
	for _, err := range errs {
		response[err.Field] = err.Message
	}

	return echo.NewHTTPError(http.StatusBadRequest, response)
}
