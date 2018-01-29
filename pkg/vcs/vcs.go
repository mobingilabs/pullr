package vcs

import (
	"context"
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

type VcsClient interface {
	CheckFileExists(ctx context.Context, repository *domain.Repository, path string, ref string) (bool, error)
	CloneRepository(ctx context.Context, repository *domain.Repository, clonePath string, ref string) error
}
