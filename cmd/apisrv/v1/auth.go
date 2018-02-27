package v1

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/auth"
)

type creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authenticatedHandlerFunc func(user string, c echo.Context) error

func (a *API) authenticated(handler authenticatedHandlerFunc) echo.HandlerFunc {
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
		return handler(user, c)
	}
}

func (a *API) login(c echo.Context) error {
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

func (a *API) register(c echo.Context) error {
	type reqPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	payload := new(reqPayload)
	if err := c.Bind(payload); err != nil {
		return err
	}

	err := a.Auth.Register(payload.Username, payload.Email, payload.Password)
	if err != nil {
		return err
	}

	secrets, err := a.Auth.Login(payload.Username, payload.Password)
	if err != nil {
		return err
	}

	setAuthSecrets(c, secrets)
	return c.NoContent(http.StatusOK)
}

func setAuthSecrets(c echo.Context, secrets *auth.Secrets) {
	c.Response().Header().Add(HeaderAuthToken, secrets.AuthToken)
	c.Response().Header().Add(HeaderRefreshToken, secrets.RefreshToken)
}
