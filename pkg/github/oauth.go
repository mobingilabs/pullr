package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// OAuthProvider implements domain.OAuthProvider interface.
type OAuthProvider struct {
	clientID     string
	clientSecret string
}

// NewOAuthProvider creates a new github oauth provider
func NewOAuthProvider(opts domain.OAuthProviderConfig) *OAuthProvider {
	return &OAuthProvider{opts.ClientID, opts.ClientSecret}
}

// LoginURL reports github login url for the user
func (g *OAuthProvider) LoginURL(secret string, cbURL string) string {
	scopes := []string{"admin:repo_hook", "read:org", "repo"}

	params := url.Values{
		"client_id":    {g.clientID},
		"scope":        {strings.Join(scopes, " ")},
		"state":        {secret},
		"redirect_uri": {cbURL},
	}.Encode()

	return fmt.Sprintf("https://github.com/login/oauth/authorize?%s", params)
}

// FinishLogin extracts necessary information from given callback request made by
// github, and finishes logging in process
func (g *OAuthProvider) FinishLogin(secret string, req *http.Request) (string, error) {
	code := req.URL.Query().Get("code")
	if strings.TrimSpace(code) == "" {
		return "", domain.ErrOAuthBadPayload
	}

	params := url.Values{
		"client_id":     {g.clientID},
		"client_secret": {g.clientSecret},
		"code":          {code},
		"state":         {secret},
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
		return "", domain.ErrOAuthBadPayload
	}

	type Payload struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}

	payload := new(Payload)
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(payload); err != nil {
		return "", domain.ErrOAuthBadPayload
	}

	if strings.ToLower(payload.TokenType) != "bearer" {
		return "", domain.ErrOAuthBadToken
	}

	return payload.AccessToken, nil
}

// Identity reports back the identity of the authenticated user as known by the github
func (*OAuthProvider) Identity(token string) (string, error) {
	return NewClient().identity(context.Background(), token)
}

// GetSecret extracts the secret from the given request
func (*OAuthProvider) GetSecret(req *http.Request) string {
	return req.URL.Query().Get("state")
}
