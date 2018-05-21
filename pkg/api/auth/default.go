package auth

import (
	"errors"
	"strings"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// DefaultAuthenticator handles authentication internally by using
// storage drivers.
type DefaultAuthenticator struct {
	authsvc *domain.DefaultAuthService
}

// NewDefaultAuthenticator creates an api authenticator
func NewDefaultAuthenticator(authsvc *domain.DefaultAuthService) *DefaultAuthenticator {
	return &DefaultAuthenticator{authsvc}
}

// Middleware creates an authentication middleware for echo server.
// Middleware authenticates the incoming request and sets proper context
// values with the user identity
func (a *DefaultAuthenticator) Middleware() echo.MiddlewareFunc {
	return func(handler echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authToken := c.Request().Header.Get(echo.HeaderAuthorization)
			authToken = strings.TrimPrefix(authToken, "Bearer ")
			if authToken == "" {
				return domain.ErrAuthUnauthorized
			}

			refreshToken := c.Request().Header.Get("X-Refresh-Token")
			if refreshToken == "" {
				return domain.ErrAuthUnauthorized
			}

			newSecrets, err := a.authsvc.Grant(refreshToken, authToken)
			if err != nil {
				return err
			}

			a.SendSecrets(c, newSecrets)
			c.Set("auth:secrets", newSecrets)
			return handler(c)
		}
	}
}

// SendSecrets appends authentication secrets to response body
func (*DefaultAuthenticator) SendSecrets(c echo.Context, secrets domain.AuthSecrets) {
	header := c.Response().Header()
	header.Set("X-Auth-Token", secrets.AuthToken)
	header.Set("X-Refresh-Token", secrets.RefreshToken)
}

// Wrap converts pullr's authenticated request handler to generic echo.HandlerFunc
func (*DefaultAuthenticator) Wrap(handler HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		secrets, ok := c.Get("auth:secrets").(domain.AuthSecrets)
		if !ok {
			return errors.New("auth wrapper couldn't find auth secrets in the context")
		}

		return handler(secrets, c)
	}
}
