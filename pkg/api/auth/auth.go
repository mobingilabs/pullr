package auth

import (
	"errors"
	"strings"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// Middleware creates an authentication middleware for echo servers. If authentication
// fails this middleware returns early and response with authentication error
func Middleware(authsvc *domain.AuthService) echo.MiddlewareFunc {
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

			newSecrets, err := authsvc.Grant(refreshToken, authToken)
			if err != nil {
				return err
			}

			SendSecrets(c, newSecrets)
			c.Set("auth:secrets", newSecrets)
			return handler(c)
		}
	}
}

// SendSecrets, writes authentication secrets to the given echo context
func SendSecrets(ctx echo.Context, secrets domain.AuthSecrets) {
	header := ctx.Response().Header()
	header.Set("X-Auth-Token", secrets.AuthToken)
	header.Set("X-Refresh-Token", secrets.RefreshToken)
}

// HandlerFunc is authenticated request handler function, takes AuthSecrets as
// it's first parameter along with echo.Context. It can be wrapped with Wrap
// function to pass it into echo router.
type HandlerFunc func(secrets domain.AuthSecrets, ctx echo.Context) error

// Wrap transforms auth.HandlerFunc into an echo.HandlerFunc
func Wrap(handler HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		secrets, ok := c.Get("auth:secrets").(domain.AuthSecrets)
		if !ok {
			return errors.New("auth wrapper couldn't find auth secrets in the context")
		}

		return handler(secrets, c)
	}
}
