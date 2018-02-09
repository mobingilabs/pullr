package v1

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/srv"
)

// OAuthLoginUrl reports OAuth authorization url for the requested oauth
// provider. Reported url includes a base64 encoded identity token (jwt) to make sure callback
// endpoint matches granted oauth token with the correct pullr user account.
func (a *apiv1) OAuthLoginUrl(username string, c echo.Context) error {
	p, ok := a.OAuthProviders[c.Param("provider")]
	if !ok {
		msg := fmt.Sprintf("Unsupported oauth provider: '%s'", c.Param("provider"))
		glog.Errorf(msg)
		return srv.NewErr("ERR_UNSUPPORTED_OAUTHPROVIDER", http.StatusBadRequest, msg)
	}

	clientUri := c.QueryParam("cb")
	clientUriTrusted := false
	for _, uri := range a.Conf.RedirectWhitelist {
		if strings.HasPrefix(clientUri, uri) {
			clientUriTrusted = true
			break
		}
	}
	if !clientUriTrusted {
		glog.Error("Untrusted uri is given for redirect, ignoring")
		return srv.NewErrBadValue("cb", clientUri)
	}

	id, err := a.Auth.NewOAuthCbIdentifier(username, p.Name(), clientUri)
	if err != nil {
		return err
	}

	cbUri := fmt.Sprintf("%s/api/v1/oauth/%s/cb/%s", a.Conf.ServerUrl, p.Name(), id.Uuid)

	loginUrl := p.LoginUrl(cbUri)
	return c.JSON(http.StatusOK, map[string]string{"login_url": loginUrl})
}

// OAuthCb handles OAuth authorization callback requests. Callback requests
// required to have an base64 encoded identity token which includes redirect url
// too. With identity token, granted OAuth token is written to correct user's
// token list. Redirect uri should start with one of the uris set by
// RedirectWhitelist configuration.
func (a *apiv1) OAuthCb(c echo.Context) (err error) {
	p, ok := a.OAuthProviders[c.Param("provider")]
	if !ok {
		msg := fmt.Sprintf("Unsupported oauth provider: '%s'", c.Param("provider"))
		return srv.NewErr("ERR_UNSUPPORTED_OAUTHPROVIDER", http.StatusBadRequest, msg)
	}

	authErr := srv.NewErr("ERR_OAUTH_FAIL", http.StatusUnauthorized, "Failed to authenticate with %s")
	errParams := errToQueryParams(authErr)

	id := c.Param("id")
	cbId, err := a.Auth.OAuthCbIdentifier(id)
	if err != nil {
		glog.Error("OAuth identifier is not provided")
		return redirect(c, a.Conf.FrontendUrl, p.Name(), errParams)
	}

	err = a.Auth.RemoveOAuthCbIdentifier(id)
	if err != nil {
		glog.Warningf("Failed to remove oauth cb identifier: %s", err)
	}

	oauthToken, err := p.HandleCb(c)
	if err != nil {
		glog.Errorln("OAuth callback couldn't handle the callback")
		return redirect(c, a.Conf.FrontendUrl, p.Name(), errParams)
	}

	err = a.Storage.PutUserToken(cbId.Username, p.Name(), oauthToken)
	if err != nil {
		glog.Errorln("OAuth callback failed to put the token into storage")
		params := errToQueryParams(srv.NewErr("ERR_INTERNAL", http.StatusInternalServerError, "Internal server error"))
		return redirect(c, cbId.RedirectUri, p.Name(), params)
	}

	return redirect(c, cbId.RedirectUri, p.Name(), url.Values{})
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
	redirectUri := appendQueryParams(uri, params)
	return c.Redirect(http.StatusTemporaryRedirect, redirectUri)
}
