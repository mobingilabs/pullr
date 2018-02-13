package vcs

import (
	"context"
	"net/http"
	"time"

	"github.com/mobingilabs/pullr/pkg/domain"
)

// RefType can be either 'branch' or 'tag'
type RefType string

// Ref types for the commit
const (
	Branch RefType = "branch"
	Tag    RefType = "tag"
)

// Webhook events
const (
	PushEvent = "push"
)

// WebhookRequest has information to work with wide range
// of VCS webhook requests.
type WebhookRequest struct {
	Event string
	Body  []byte
}

// CommitInfo has information about a commit
type CommitInfo struct {
	// Author of the commit
	Author string
	// Ref is git tag for the tagged commits otherwise the branch name
	Ref string
	// RefType says what ref value is about either Branch or Tag
	RefType RefType
	// Hash is the commit id hash
	Hash string
	// CreatedAt is time of the commit
	CreatedAt time.Time
	// Repository is the source code repository
	Repository domain.Repository
}

// Vcs abstracts anonymous vcs operations
type Vcs interface {
	// ExtractCommitInfo parses the given WebhookRequest and reports back the
	// commit info.
	ExtractCommitInfo(r *WebhookRequest) (*CommitInfo, error)

	// ParseWebhookRequest parses given http.Request as a WebhookRequest.
	ParseWebhookRequest(r *http.Request) (*WebhookRequest, error)
}

// Client abstracts authenticated vcs client operations
type Client interface {
	// ListOrganisations fetches authenticated user's organisations
	ListOrganisations(ctx context.Context) ([]string, error)

	// CheckFileExists checks a repository if the given file path exists or not
	CheckFileExists(ctx context.Context, repository domain.Repository, path string, ref string) (bool, error)

	// CloneRepository clones the repository content to given path
	CloneRepository(ctx context.Context, repository domain.Repository, clonePath string, ref string) error

	// ListRepositories fetches authenticated user's repositories
	ListRepositories(ctx context.Context, organisation string) ([]string, error)
}
