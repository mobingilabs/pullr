package domain

import (
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
	ParseWebhookPayload(req *http.Request) (*CommitInfo, error)
	Organisations(username string, token string) ([]string, error)
	Repositories(username string, organisation string, token string) ([]string, error)
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
