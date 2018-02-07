package v1

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

var ErrOauthBadProvider = errors.New("bad provider")

func (a *apiv1) OAuthLoginUrl(username string, c echo.Context) error {
	p := c.Param("provider")
	if !isOAuthProviderValid(p) {
		return ErrOauthBadProvider
	}

	cbUri := c.QueryParam("cb")
	if cbUri == "" {
		return c.JSON(http.StatusBadRequest, struct {
			Message string `json:"message"`
		}{"Missing parameter 'cb'"})
	}

	cbUri64 := base64.StdEncoding.EncodeToString([]byte(cbUri))

	// OAuth callback url http[s]://SERVER_URL/oauth/PROVIDER/cb/FRONTEND_URL_BASE64
	authRedirUrl := fmt.Sprintf("%s/api/v1/oauth/%s/cb/%s", a.Conf.ServerUrl, p, cbUri64)

	loginUrl := ""
	switch p {
	case "github":
		scope := ""
		loginUrl = fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", a.Conf.GithubClientId, authRedirUrl)
	default:
		return ErrOauthBadProvider
	}

	return c.JSON(http.StatusOK, map[string]string{"login_url": loginUrl})
}

func (a *apiv1) OAuthPutToken(username string, c echo.Context) error {
	provider := c.Param("provider")
	if !isOAuthProviderValid(provider) {
		return ErrOauthBadProvider
	}

	type Payload struct {
		Token string `json:"token"`
	}

	payload := new(Payload)
	if err := c.Bind(payload); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	a.Storage.PutUserToken(username, provider, payload.Token)
	return c.NoContent(http.StatusCreated)
}

func (a *apiv1) OAuthCb(c echo.Context) error {
	provider := c.Param("provider")
	if !isOAuthProviderValid(provider) {
		// TODO: Better to redirect user to client provided redirect uri?
		return c.NoContent(http.StatusBadRequest)
	}

	redirectUri := c.Param("redirectUri")
	if redirectUri == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	return a.handleOAuthCb(provider, c)
}

func (a *apiv1) handleOAuthCb(provider string, c echo.Context) error {
	switch provider {
	case "github":
		code := c.QueryParam("code")
	}
}

func isOAuthProviderValid(provider string) bool {
	switch provider {
	case "github":
		return true
	default:
		return false
	}
}
