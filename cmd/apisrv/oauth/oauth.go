package oauth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// Errors
var (
	ErrUnexpected        = errors.New("authentication failed")
	ErrUnsupportedToken  = errors.New("unsupported token type")
	ErrUnexpectedPayload = errors.New("unexpected payload")
)

type Perm int

type Client interface {
	// Name reports OAuthProvider's name
	Name() string

	// LoginUrl reports a valid oauth login url for the client to visit
	LoginUrl(cbUrl string) string

	// HandleCb handles callback request made by OAuth provider and reports back
	// the access token.
	HandleCb(c echo.Context) (string, error)
}

type CbClaims struct {
	jwt.StandardClaims
	RedirectUri string `json:"redirect_uri"`
}
