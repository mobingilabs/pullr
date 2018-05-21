package auth

import (
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// HandlerFunc is authenticated request handler function, takes AuthSecrets as
// it's first parameter along with echo.Context. It can be wrapped with Wrap
// function to pass it into echo router.
type HandlerFunc func(secrets domain.AuthSecrets, ctx echo.Context) error

// Authenticator authenticate incoming requests
type Authenticator interface {
	// Middleware creates an authentication middleware for echo server.
	// Middleware authenticates the incoming request and sets proper context
	// values with the user identity
	Middleware() echo.MiddlewareFunc

	// SendSecrets appends authentication secrets to response body
	SendSecrets(c echo.Context, secrets domain.AuthSecrets)

	// Wrap converts pullr's authenticated request handler to generic echo.HandlerFunc
	Wrap(handler HandlerFunc) echo.HandlerFunc
}
