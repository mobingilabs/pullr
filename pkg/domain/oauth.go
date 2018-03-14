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
}

// OAuthStorage handles storing and querying oauth related data
type OAuthStorage interface {
	// PutSecret puts a new oauth login secret for making sure further
	// incoming oauth callback requests are really made by the provider
	PutSecret(username, secret string) error

	// PopSecret removes a secret from the storage, if the given secret
	// is not found it returns notfound error.
	PopSecret(secret string) (string, error)

	// GetTokens finds matching oauth tokens for given user
	GetTokens(username string) (map[string]string, error)

	// PutToken inserts a new token record for the given user
	PutToken(username string, provider string, token string) error

	// RemoveToken removes a token from the given user
	RemoveToken(username string, provider string) error
}

// OAuthProvider provides helpers for logging with a specific oauth provider
type OAuthProvider interface {
	// LoginUrl reports an oauth provider login url to redirect users.
	// secret value should be persisted til callback request handled
	LoginUrl(secret string, cbUrl string) string

	// HandleCallback handles oauth provider's incoming callback request
	// to finish logging in process. It reports back the received token.
	HandleCallback(secret string, req *http.Request) (string, error)

	// GetSecret reports the secret given by the oauth provider from the
	// callback request
	GetSecret(req *http.Request) string
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

// LoginUrl reports back a login url to oauth provider's login endpoint
func (s *OAuthService) LoginUrl(provider string, username string, cbUrl string) (string, error) {
	p, ok := s.providers[provider]
	if !ok {
		return "", ErrOAuthUnsupportedProvider
	}

	randBytes := make([]byte, 16)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}

	secret := base64.StdEncoding.EncodeToString(randBytes)

	err = s.storage.PutSecret(username, secret)
	if err != nil {
		return "", err
	}

	return p.LoginUrl(secret, cbUrl), nil
}

// HandleCallback processes received callback request from the oauth provider
// and finalizes authentication by getting a token from the provider
func (s *OAuthService) HandleCallback(provider string, req *http.Request) (string, error) {
	p, ok := s.providers[provider]
	if !ok {
		return "", ErrOAuthUnsupportedProvider
	}

	secret := p.GetSecret(req)
	username, err := s.storage.PopSecret(secret)
	if err != nil {
		return "", err
	}
	if username == "" {
		return "", ErrNotFound
	}

	return p.HandleCallback(secret, req)
}
