package v1

import (
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// UserProfile is an handler serves only authenticated users. Responds with
// authenticated user's profile data
func (a *Api) UserProfile(secrets domain.AuthSecrets, c echo.Context) error {
	return nil
}

// UserProfileUpdate is an handle
func (a *Api) UserProfileUpdate(secrets domain.AuthSecrets, c echo.Context) error {
	return nil
}
