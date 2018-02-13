package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/oauth"
	"github.com/mobingilabs/pullr/pkg/srv"
)

// github is an OAuth provider
type github struct {
	clientID string
	secret   string
}

// New creates a github specific oauth client
func New(clientID, secret string) oauth.Client {
	return &github{clientID, secret}
}

func (g *github) Name() string {
	return "github"
}

func (g *github) LoginURL(cbURL string) string {
	scopes := []string{"admin:repo_hook", "read:org", "repo"}

	params := url.Values{
		"client_id":    {g.clientID},
		"scope":        {strings.Join(scopes, " ")},
		"redirect_uri": {cbURL},
	}.Encode()

	return fmt.Sprintf("https://github.com/login/oauth/authorize?%s", params)
}

func (g *github) HandleCb(r *http.Request) (string, error) {
	code := r.URL.Query().Get("code")
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
		return "", oauth.ErrUnexpected
	}

	type Payload struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}

	payload := new(Payload)
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(payload); err != nil {
		return "", oauth.ErrUnexpectedPayload
	}

	if strings.ToLower(payload.TokenType) != "bearer" {
		return "", oauth.ErrUnsupportedToken
	}

	return payload.AccessToken, nil
}
