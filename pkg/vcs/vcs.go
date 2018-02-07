package vcs

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/mobingilabs/pullr/pkg/domain"
)

// Ref types for the commit
const (
	Branch = "branch"
	Tag    = "tag"
)

// Webhook events
const (
	PushEvent = "push"
)

type OAuthPerm int

// OAuth permissions
const (
	PermReadRepos OAuthPerm = iota
	PermListOrgs
	PermAdminRepoHooks
)

var (
	ErrAuthRequired   = errors.New("vcs client needs to authenticate")
	ErrOAuthInvalidCb = errors.New("invalid oauth cb url provided")
)

type OAuthSecrets struct {
	ClientId     string
	ClientSecret string
}

type WebhookRequest struct {
	Event string
	Body  []byte
}

type CommitInfo struct {
	Author     string
	Ref        string
	RefType    string
	Hash       string
	CreatedAt  time.Time
	Repository domain.Repository
}

type Vcs interface {
	ExtractCommitInfo(r *WebhookRequest) (*CommitInfo, error)
	ParseWebhookRequest(r *http.Request) (*WebhookRequest, error)
}

type Client interface {
	ListOrganisations(ctx context.Context) ([]string, error)
	CheckFileExists(ctx context.Context, repository domain.Repository, path string, ref string) (bool, error)
	CloneRepository(ctx context.Context, repository domain.Repository, clonePath string, ref string) ([]byte, error)
	ListRepositories(ctx context.Context, organisation string) ([]string, error)
}
