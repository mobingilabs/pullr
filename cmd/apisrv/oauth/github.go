package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/srv"
)

// Github is an OAuth provider
type Github struct {
	clientID string
	secret   string
}

// NewGithub creates a new github oauth implementation instance
func NewGithub(clientID, secret string) *Github {
	return &Github{clientID, secret}
}

// Name reports oauth provider's name
func (g *Github) Name() string {
	return "github"
}

// LoginURL reports login url for the oauth provider instance
func (g *Github) LoginURL(cbURL string) string {
	scopes := []string{"admin:repo_hook", "read:org", "repo"}

	params := url.Values{
		"client_id":    {g.clientID},
		"scope":        {strings.Join(scopes, " ")},
		"redirect_uri": {cbURL},
	}.Encode()

	return fmt.Sprintf("https://github.com/login/oauth/authorize?%s", params)
}

// HandleCb is handles github's callback request to fetch oauth token
func (g *Github) HandleCb(c echo.Context) (string, error) {
	code := c.QueryParam("code")
	if code == "" || strings.TrimSpace(code) == "" {
		return "", srv.NewErrMissingParam("code")
	}
	params := url.Values{
		"client_id":     {g.clientID},
		"client_secret": {g.secret},
		"code":          {code},
	}.Encode()

	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?%s", params)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
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
