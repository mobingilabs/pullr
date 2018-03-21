package v1

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// UserProfile is an handler serves only authenticated users. Responds with
// authenticated user's profile data
func (a *Api) UserProfile(secrets domain.AuthSecrets, c echo.Context) error {
	usr, err := a.userStorage.Get(secrets.Username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

// UserProfileUpdate is an handle
func (a *Api) UserProfileUpdate(secrets domain.AuthSecrets, c echo.Context) error {
	return nil
}
