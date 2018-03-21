package github

import (
	"time"

	"github.com/mobingilabs/pullr/pkg/gova"
)

// PushEvent represents a git push to a GitHub repository.
//
// Actually push events contains more data than described here. This definition
// only contains pullr related fields to keep it simple.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#pushevent
type PushEvent struct {
	Ref        *string `json:"ref"`
	Before     *string `json:"before"`
	After      *string `json:"after"`
	HeadCommit *struct {
		Author *struct {
			Name *string `json:"name"`
		} `json:"author"`

		Timestamp *time.Time `json:"timestamp,omitempty"`
	} `json:"head_commit,omitempty"`

	Repository *struct {
		Name  *string `json:"name"`
		Owner *struct {
			Login *string `json:"login"`
		} `json:"owner,omitempty"`
	} `json:"repository,omitempty"`
}

// Validate validates the push event
func (p *PushEvent) Validate() (bool, error) {
	val := &gova.Validator{}
	val.NotNil("ref", p.Ref)
	val.NotNil("before", p.Before)
	val.NotNil("before", p.After)
	val.NotNil("head_commit", p.HeadCommit)
	if p.HeadCommit != nil {
		val.NotNil("head_commit.author", p.HeadCommit.Author)
		if p.HeadCommit.Author != nil {
			val.NotNil("head_commit.author.name", p.HeadCommit.Author.Name)
		}

		val.NotNil("head_commit.timestamp", p.HeadCommit.Timestamp)
	}

	val.NotNil("repository", p.Repository)
	if p.Repository != nil {
		val.NotNil("repository.name", p.Repository.Name)
		val.NotNil("repository.owner", p.Repository.Owner)
		if p.Repository.Owner != nil {
			val.NotNil("repository.owner.login", p.Repository.Owner.Login)
		}
	}

	return val.Valid(), val.Errors()
}
