package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/cmd/apisrv/perrors"
)

type Github struct {
	clientId string
	secret   string
}

func NewGithub(clientId, secret string) *Github {
	return &Github{clientId, secret}
}

func (g *Github) Name() string {
	return "github"
}

func (g *Github) LoginUrl(cbUrl string) string {
	scopes := []string{"admin:repo_hook", "read:org", "repo"}

	params := url.Values{
		"client_id":    {g.clientId},
		"scope":        {strings.Join(scopes, " ")},
		"redirect_uri": {cbUrl},
	}.Encode()

	return fmt.Sprintf("https://github.com/login/oauth/authorize?%s", params)
}

func (g *Github) HandleCb(c echo.Context) (string, error) {
	code := c.QueryParam("code")
	if code == "" || strings.TrimSpace(code) == "" {
		return "", perrors.NewErrMissingParam("code")
	}
	params := url.Values{
		"client_id":     {g.clientId},
		"client_secret": {g.secret},
		"code":          {code},
	}.Encode()

	reqUrl := fmt.Sprintf("https://github.com/login/oauth/access_token?%s", params)
	req, err := http.NewRequest(http.MethodPost, reqUrl, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add(echo.HeaderAccept, echo.MIMEApplicationJSON)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	successCode := res.StatusCode >= 200 && res.StatusCode < 300
	if !successCode {
		return "", ErrUnexpected
	}

	type Payload struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}

	payload := new(Payload)
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(payload); err != nil {
		return "", ErrUnexpectedPayload
	}

	if strings.ToLower(payload.TokenType) != "bearer" {
		return "", ErrUnsupportedToken
	}

	return payload.AccessToken, nil
}
