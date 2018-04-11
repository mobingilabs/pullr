package domain

import (
	"context"
	"net/http"
	"time"
)

// SourceRefType can be either 'branch' or 'tag'
type SourceRefType string

// Ref types for the commit
const (
	SourceBranch SourceRefType = "branch"
	SourceTag    SourceRefType = "tag"
)

// SourceClient wraps source control provider client operations
type SourceClient interface {
	// ParseWebhookPayload extracts CommitInfo from source provider's webhook request
	ParseWebhookPayload(req *http.Request) (*CommitInfo, error)

	// RegisterWebhook registers pullr to source provider's webhooks
	RegisterWebhook(ctx context.Context, token string, webhookURL string, repo SourceRepository) error

	// Organisations reports back a list of organisations of the authenticated source
	// provider user
	Organisations(ctx context.Context, identity string, token string) ([]string, error)

	// Repositories reports back a list of repositories which belongs to given
	// organisation and user can collaborate.
	//
	// To get repositories owned by the user pass authenticated source provider
	// user name as the organisation.
	Repositories(ctx context.Context, identity string, organisation string, token string) ([]string, error)
}

// CommitInfo has information about a commit
type CommitInfo struct {
	// Author of the commit
	Author string
	// Ref is git tag for the tagged commits otherwise the branch name
	Ref string
	// RefType says what ref value is about either Branch or Tag
	RefType SourceRefType
	// Hash is the commit id hash
	Hash string
	// CreatedAt is time of the commit
	CreatedAt time.Time
	// SourceRepository is the source code repository
	Repository SourceRepository
}

// SourceRepository has the information for source code repository.
type SourceRepository struct {
	Provider string `json:"provider" bson:"provider"`
	Owner    string `json:"owner" bson:"owner"`
	Name     string `json:"name" bson:"name"`
}

// SourceService wraps different vcs client implementations
// to integrate with vcs providers.
type SourceService struct {
	storage OAuthStorage
	clients map[string]SourceClient
}

// NewSourceService creates a new SourceService
func NewSourceService(storage OAuthStorage, clients map[string]SourceClient) *SourceService {
	return &SourceService{storage, clients}
}

// ParseWebhookPayload extracts the commit info from given source provider's webhook request.
func (s *SourceService) ParseWebhookPayload(provider string, req *http.Request) (*CommitInfo, error) {
	c, ok := s.clients[provider]
	if !ok {
		return nil, ErrSourceUnsupportedProvider
	}

	return c.ParseWebhookPayload(req)
}

// RegisterWebhook registers pullr to source provider's webhooks
func (s *SourceService) RegisterWebhook(ctx context.Context, webhookURL, username string, repo SourceRepository) error {
	c, ok := s.clients[repo.Provider]
	if !ok {
		return ErrSourceUnsupportedProvider
	}

	tokens, err := s.storage.GetTokens(username)
	if err != nil {
		return err
	}

	ptoken, ok := tokens[repo.Provider]
	if !ok {
		return ErrAuthUnauthorized
	}

	return c.RegisterWebhook(ctx, ptoken.Token, webhookURL, repo)
}

// Organisations find organisations which user has membership
func (s *SourceService) Organisations(ctx context.Context, provider, username string) ([]string, error) {
	c, ok := s.clients[provider]
	if !ok {
		return nil, ErrSourceUnsupportedProvider
	}

	tokens, err := s.storage.GetTokens(username)
	if err != nil {
		return nil, err
	}

	token, ok := tokens[provider]
	if !ok {
		return nil, ErrAuthUnauthorized
	}

	return c.Organisations(ctx, token.Identity, token.Token)
}

// Repositories finds repositories belongs to organisation
func (s *SourceService) Repositories(ctx context.Context, provider, username, organisation string) ([]string, error) {
	c, ok := s.clients[provider]
	if !ok {
		return nil, ErrSourceUnsupportedProvider
	}

	tokens, err := s.storage.GetTokens(username)
	if err != nil {
		return nil, err
	}

	token, ok := tokens[provider]
	if !ok {
		return nil, ErrAuthUnauthorized
	}

	return c.Repositories(ctx, token.Identity, organisation, token.Token)
}
