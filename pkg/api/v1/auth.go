package v1

import (
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/api/auth"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// AuthLogin is an handler for login requests. Authenticates the user
// if their credentials are correct.
func (a *Api) AuthLogin(c echo.Context) error {
	type loginPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var payload loginPayload
	if err := c.Bind(&payload); err != nil {
		return err
	}

	secrets, err := a.authsvc.Login(payload.Username, payload.Password)
	if err != nil {
		return err
	}

	auth.SendSecrets(c, secrets)
	return nil
}

// AuthRegister is an handler for register requests. It stores both
// user data and their credentials, and then authenticates the user.
func (a *Api) AuthRegister(c echo.Context) error {
	type registerPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	var payload registerPayload
	if err := c.Bind(&payload); err != nil {
		return err
	}

	err := a.authsvc.Register(payload.Username, payload.Email, payload.Password)
	if err != nil {
		return err
	}

	err = a.userStorage.Put(domain.User{
		Username: payload.Username,
		Email:    payload.Email,
	})
	if err != nil {
		a.authStorage.Delete(payload.Username)
		return err
	}

	secrets, err := a.authsvc.Login(payload.Username, payload.Password)
	if err != nil {
		return err
	}

	auth.SendSecrets(c, secrets)
	return nil
}
