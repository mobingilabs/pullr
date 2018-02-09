package v1

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/auth"
	"github.com/mobingilabs/pullr/pkg/srv"
)

type creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authenticatedHandlerFunc func(user string, c echo.Context) error
type tokenKind int

const (
	HeaderAuthToken    = "X-Auth-Token"
	HeaderRefreshToken = "X-Refresh-Token"

	TokenAuth tokenKind = iota
	TokenRefresh
)

func (a *apiv1) authenticated(h authenticatedHandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authToken := c.Request().Header.Get(echo.HeaderAuthorization)
		authToken = strings.TrimPrefix(authToken, "Bearer ")
		if authToken == "" {
			return auth.ErrCredentials
		}

		refreshToken := c.Request().Header.Get(HeaderRefreshToken)
		if refreshToken == "" {
			return auth.ErrCredentials
		}

		newSecrets, user, err := a.Auth.Validate(authToken, refreshToken)
		if err != nil {
			return err
		}

		setAuthSecrets(c, newSecrets)
		return h(user, c)
	}
}

func (a *apiv1) getToken(c echo.Context, kind tokenKind) (*jwt.Token, *jwt.StandardClaims, error) {
	switch kind {
	case TokenAuth:
		bearer := c.Request().Header.Get(echo.HeaderAuthorization)
		bearer = strings.TrimPrefix(bearer, "Bearer ")
		claims := new(jwt.StandardClaims)
		token, err := a.Auth.ParseToken(bearer, claims)
		return token, claims, err
	case TokenRefresh:
		tokenStr := c.Request().Header.Get(HeaderRefreshToken)
		claims := new(jwt.StandardClaims)
		token, err := a.Auth.ParseToken(tokenStr, claims)
		return token, claims, err
	}

	return nil, nil, srv.NewErr("ERR_INTERNAL", http.StatusInternalServerError, "An unexpected error happened")
}

func (a *apiv1) login(c echo.Context) error {
	credentials := new(creds)
	if err := c.Bind(credentials); err != nil {
		return err
	}

	secrets, err := a.Auth.Login(credentials.Username, credentials.Password)
	if err != nil {
		return err
	}

	setAuthSecrets(c, secrets)
	return c.NoContent(http.StatusOK)
}

func (a *apiv1) register(c echo.Context) error {
	credentials := new(creds)
	if err := c.Bind(credentials); err != nil {
		return err
	}

	err := a.Auth.Register(credentials.Username, credentials.Password)
	if err != nil {
		return err
	}

	secrets, err := a.Auth.Login(credentials.Username, credentials.Password)
	setAuthSecrets(c, secrets)

	return c.NoContent(http.StatusOK)
}

func setAuthSecrets(c echo.Context, secrets *auth.Secrets) {
	c.Response().Header().Add(HeaderAuthToken, secrets.AuthToken)
	c.Response().Header().Add(HeaderRefreshToken, secrets.RefreshToken)
}
