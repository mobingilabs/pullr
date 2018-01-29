package v1

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/auth"
)

type creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authenticatedHandlerFunc func(user string, c echo.Context) error

const (
	authCookieName    = "AuthToken"
	refreshCookieName = "RefreshToken"
	csrfHeaderName    = "X-CSRF-TOKEN"
)

func (a *apiv1) authenticated(h authenticatedHandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authCookie, err := c.Cookie(authCookieName)
		if err != nil {
			if err != http.ErrNoCookie {
				return err
			}

			return auth.ErrCredentials
		}

		refreshCookie, err := c.Cookie(refreshCookieName)
		if err != nil {
			if err != http.ErrNoCookie {
				return err
			}

			return auth.ErrCredentials
		}

		csrf := c.Request().Header.Get(csrfHeaderName)
		if csrf == "" {
			return auth.ErrCredentials
		}

		newSecrets, user, err := a.Auth.Validate(csrf, authCookie.Value, refreshCookie.Value)
		if err != nil {
			return err
		}

		setAuthSecrets(c, newSecrets)
		return h(user, c)
	}
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

func (a *apiv1) logout(c echo.Context) error {
	nullifyAuthSecrets(c)
	return c.NoContent(http.StatusOK)
}

func setAuthSecrets(c echo.Context, secrets *auth.Secrets) {
	authCookie := &http.Cookie{
		Name:     authCookieName,
		Value:    secrets.AuthToken,
		HttpOnly: true,
	}

	refreshCookie := &http.Cookie{
		Name:     refreshCookieName,
		Value:    secrets.RefreshToken,
		HttpOnly: true,
	}

	c.SetCookie(authCookie)
	c.SetCookie(refreshCookie)
	c.Response().Header().Add(csrfHeaderName, secrets.Csrf)
}

func nullifyAuthSecrets(c echo.Context) {
	authCookie := &http.Cookie{
		Name:     authCookieName,
		Value:    "",
		Expires:  time.Now().Add(time.Hour * -1000),
		HttpOnly: true,
	}

	refreshToken := &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Expires:  time.Now().Add(time.Hour * -1000),
		HttpOnly: true,
	}

	c.SetCookie(authCookie)
	c.SetCookie(refreshToken)
	c.Response().Header().Add(csrfHeaderName, "")
}
