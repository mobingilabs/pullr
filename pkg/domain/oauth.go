package domain

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

// OAuthToken represents a token received from an oauth provider
type OAuthToken struct {
	Provider string `json:"provider"`
	Token    string `json:"token"`
	Identity string `json:"identity"`
	Redirect string `json:"redir"`
}

// OAuthStorage handles storing and querying oauth related data
type OAuthStorage interface {
	// PutSecret puts a new oauth login secret for making sure further
	// incoming oauth callback requests are really made by the provider
	PutSecret(username, secret, cburi string) error

	// PopSecret returns secret associated callback url, if the given
	// secret is not found it returns notfound error.
	PopSecret(username, secret string) (string, error)

	// GetTokens finds matching oauth tokens for given user
	GetTokens(username string) (map[string]OAuthToken, error)

	// PutToken inserts a new token record for the given user
	PutToken(username string, identity string, provider string, token string) error

	// RemoveToken removes a token from the given user
	RemoveToken(username string, provider string) error
}

// OAuthProvider provides helpers for logging with a specific oauth provider
type OAuthProvider interface {
	// LoginURL reports an oauth provider login url to redirect users.
	// secret value should be persisted til callback request handled
	LoginURL(secret string, cbUrl string) string

	// FinishLogin handles oauth provider's incoming callback request
	// to finish logging in process. It reports back the received token.
	FinishLogin(secret string, req *http.Request) (string, error)

	// GetSecret reports the secret given by the oauth provider from the
	// callback request
	GetSecret(req *http.Request) string

	// Identity reports back the identity of the authenticated user
	// as known by the provider
	Identity(token string) (string, error)
}

// OAuthService handles oauth logging in operations with given set of oauth
// oauth provider implementations
type OAuthService struct {
	storage   OAuthStorage
	providers map[string]OAuthProvider
}

// NewOAuthService creates an oauth service to handle 3rd party oauth authentications
func NewOAuthService(storage OAuthStorage, providers map[string]OAuthProvider) *OAuthService {
	return &OAuthService{storage, providers}
}

// LoginURL reports back a login url to oauth provider's login endpoint
func (s *OAuthService) LoginURL(provider string, username string, cbUrl, redir string) (string, error) {
	p, ok := s.providers[provider]
	if !ok {
		return "", ErrOAuthUnsupportedProvider
	}

	randBytes := make([]byte, 8)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}

	secret := base64.StdEncoding.EncodeToString(randBytes)

	err = s.storage.PutSecret(username, secret, redir)
	if err != nil {
		return "", err
	}

	return p.LoginURL(secret, cbUrl), nil
}

// FinishLogin processes received callback request from the oauth provider
// and finalizes authentication by getting a token from the provider
func (s *OAuthService) FinishLogin(provider string, reqUsername string, callackReq *http.Request) (OAuthToken, error) {
	p, ok := s.providers[provider]
	if !ok {
		return OAuthToken{}, ErrOAuthUnsupportedProvider
	}

	secret := p.GetSecret(callackReq)
	redir, err := s.storage.PopSecret(reqUsername, secret)
	if err != nil {
		return OAuthToken{}, err
	}

	token, err := p.FinishLogin(secret, callackReq)
	if err != nil {
		return OAuthToken{}, err
	}

	identity, err := p.Identity(token)
	if err != nil {
		return OAuthToken{}, err
	}

	err = s.storage.PutToken(reqUsername, identity, provider, token)
	oauthToken := OAuthToken{
		Provider: provider,
		Token:    token,
		Redirect: redir,
		Identity: identity,
	}
	return oauthToken, err
}
