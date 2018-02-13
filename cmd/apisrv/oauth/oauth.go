package oauth

import (
	"errors"

	"github.com/labstack/echo"
)

// Errors
var (
	ErrUnexpected        = errors.New("authentication failed")
	ErrUnsupportedToken  = errors.New("unsupported token type")
	ErrUnexpectedPayload = errors.New("unexpected payload")
)

// Client is an OAuth provider client. Implementors are responsible for
// generating proper login urls for the browser as well as handling login
// callback requests made by the provider
type Client interface {
	// Name reports OAuthProvider's name
	Name() string

	// LoginURL reports a valid oauth login url for the client to visit
	LoginURL(cbURL string) string

	// HandleCb handles callback request made by OAuth provider and reports back
	// the access token.
	HandleCb(c echo.Context) (string, error)
}
