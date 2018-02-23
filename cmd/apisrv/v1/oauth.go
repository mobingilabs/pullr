package v1

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/srv"
	log "github.com/sirupsen/logrus"
)

// oauthLoginURL reports OAuth authorization url for the requested oauth
// provider. Reported url includes a base64 encoded identity token (jwt) to make sure callback
// endpoint matches granted oauth token with the correct pullr user account.
func (a *API) oauthLoginURL(username string, c echo.Context) error {
	p, ok := a.OAuth[c.Param("provider")]
	if !ok {
		msg := fmt.Sprintf("Unsupported oauth provider: '%s'", c.Param("provider"))
		return srv.NewErr("ERR_UNSUPPORTED_OAUTHPROVIDER", http.StatusBadRequest, msg)
	}

	clientURI := c.QueryParam("cb")
	id, err := a.Auth.NewOAuthCbIdentifier(username, p.Name(), clientURI)
	if err != nil {
		return err
	}

	cbURI := fmt.Sprintf("https://%s/api/v1/oauth/%s/cb/%s", c.Request().Host, p.Name(), id.UUID)

	loginURL := p.LoginURL(cbURI)
	return c.JSON(http.StatusOK, map[string]string{"login_url": loginURL})
}

// oauthCb handles OAuth authorization callback requests. Callback requests
// required to have an base64 encoded identity token which includes redirect url
// too. With identity token, granted OAuth token is written to correct user's
// token list. Redirect uri should start with one of the uris set by
// RedirectWhitelist configuration.
func (a *API) oauthCb(c echo.Context) (err error) {
	client, ok := a.OAuth[c.Param("provider")]
	if !ok {
		msg := fmt.Sprintf("Unsupported oauth provider: '%s'", c.Param("provider"))
		return srv.NewErr("ERR_UNSUPPORTED_OAUTHPROVIDER", http.StatusBadRequest, msg)
	}

	authErr := srv.NewErr("ERR_OAUTH_FAIL", http.StatusUnauthorized, "Failed to authenticate with %s")
	errParams := errToQueryParams(authErr)

	redirectURL := fmt.Sprintf("http://%s", c.Request().Host)

	id := c.Param("id")
	cbID, err := a.Auth.OAuthCbIdentifier(id)
	if err != nil {
		log.Warn("OAuth identifier is not provided")
		return redirect(c, redirectURL, client.Name(), errParams)
	}

	err = a.Auth.RemoveOAuthCbIdentifier(id)
	if err != nil {
		log.Warnf("Failed to remove oauth cb identifier: %s", err)
	}

	oauthToken, err := client.HandleCb(c.Request())
	if err != nil {
		log.Warn("OAuth callback couldn't handle the callback")
		return redirect(c, redirectURL, client.Name(), errParams)
	}

	err = a.Storage.PutUserToken(cbID.Username, client.Name(), oauthToken)
	if err != nil {
		log.Error("OAuth callback failed to put the token into storage")
		params := errToQueryParams(srv.NewErr("ERR_INTERNAL", http.StatusInternalServerError, "Internal server error"))
		return redirect(c, cbID.RedirectURI, client.Name(), params)
	}

	return redirect(c, cbID.RedirectURI, client.Name(), url.Values{})
}

func appendQueryParams(uri string, params url.Values) string {
	query := params.Encode()
	queryPrefix := "?"
	if strings.Contains(uri, "?") {
		queryPrefix = "&"
	}

	return fmt.Sprintf("%s%s%s", uri, queryPrefix, query)
}

func errToQueryParams(err srv.ErrMsg) url.Values {
	return url.Values{
		"err_kind":   {err.Kind},
		"err_status": {strconv.FormatInt(int64(err.Status), 10)},
	}
}

func redirect(c echo.Context, uri, provider string, params url.Values) error {
	params.Add("provider", provider)
	redirectURI := appendQueryParams(uri, params)
	return c.Redirect(http.StatusTemporaryRedirect, redirectURI)
}
