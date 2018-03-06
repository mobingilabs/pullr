package srv

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/auth"
	"github.com/mobingilabs/pullr/pkg/storage"
	"github.com/mobingilabs/pullr/pkg/vcs"
)

// ElapsedMiddleware logs elapsed time while handling the request
func ElapsedMiddleware(logger Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()

			res := next(c)

			elapsed := time.Since(startTime)
			logger.Infof("%s %s %d took %s", c.Request().Method, c.Request().URL.Path, c.Response().Status, elapsed)

			return res
		}
	}
}

// ServerHeaderMiddleware adds server name to response headers
func ServerHeaderMiddleware(srvName string, version string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderServer, "mobingi:pullr:apiserver:"+version)
			return next(c)
		}
	}
}

// ErrorMiddleware is an echo middleware to map few known error values from common
// packages as well as ErrMsg values.
func ErrorMiddleware(logger Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			logger.Error(err)

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
			case auth.ErrEmailTaken:
				e.Kind = "ERR_EMAILTAKEN"
				e.Status = http.StatusConflict
				e.Msg = "Email is already registered"

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
}
