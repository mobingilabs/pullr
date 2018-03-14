package v1

import (
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// OAuthLogin responses with oauth provider's login url with correct callback url
func (a *Api) OAuthLogin(secrets domain.AuthSecrets, c echo.Context) error {
	return nil
}

// OAuthCallback handles callback request made by an oauth provider to authorize
// Pullr to user's resources on the provider
func (a *Api) OAuthCallback(c echo.Context) error {
	return nil
}
