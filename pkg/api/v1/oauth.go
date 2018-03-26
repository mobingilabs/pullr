package v1

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// OAuthLogin responses with oauth provider's login url with correct callback url
func (a *Api) OAuthLogin(secrets domain.AuthSecrets, c echo.Context) error {
	provider := c.Param("provider")
	cburl := fmt.Sprintf("https://%s/api/v1/oauth/%s/cb/%s", c.Request().Host, provider, secrets.Username)
	loginurl, err := a.oauthsvc.LoginURL(provider, secrets.Username, cburl, c.QueryParam("redirect"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"login_url": loginurl})
}

// OAuthCallback handles callback request made by an oauth provider to authorize
// Pullr to user's resources on the provider
func (a *Api) OAuthCallback(c echo.Context) error {
	provider := c.Param("provider")
	username := c.Param("username")

	_, redir, err := a.oauthsvc.FinishLogin(provider, username, c.Request())
	if err != nil {
		return err
	}

	params := url.Values{"provider": []string{provider}}.Encode()
	sep := "?"
	if strings.Contains(redir, "?") {
		sep = "&"
	}

	redir = fmt.Sprintf("%s%s%s", redir, sep, params)
	return c.Redirect(http.StatusTemporaryRedirect, redir)
}
