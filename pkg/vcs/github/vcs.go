package github

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/vcs"
)

// Github provides Vcs functionality
type Github struct{}

func New() *Github {
	return &Github{}
}

// ExtractCommitInfo read webhook request payload and tries to extract commit info out of it
func (*Github) ExtractCommitInfo(r *vcs.WebhookRequest) (*vcs.CommitInfo, error) {
	switch r.Event {
	case vcs.PushEvent:
		return extractCommitInfoPushPayload(r)
	default:
		return nil, vcs.ErrUnsupportedEvent
	}
}

// ParseWebhookRequest tries to parse given webhook http request
func (*Github) ParseWebhookRequest(r *http.Request) (*vcs.WebhookRequest, error) {
	if !strings.HasPrefix(r.Header.Get("User-Agent"), "GitHub-Hookshot") {
		return nil, vcs.ErrInvalidWebhook
	}

	event := r.Header.Get("X-Github-Event")
	if event == "" {
		return nil, vcs.ErrInvalidWebhook
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}

	webhookRequest := &vcs.WebhookRequest{
		Event: event,
		Body:  body,
	}

	return webhookRequest, nil
}

func extractCommitInfoPushPayload(r *vcs.WebhookRequest) (*vcs.CommitInfo, error) {
	var event github.PushEvent
	if err := json.Unmarshal(r.Body, &event); err != nil {
		return nil, vcs.ErrInvalidWebhookPayload
	}

	// Parse push type
	parts := strings.Split(event.GetRef(), "/")
	refType := ""
	refName := ""
	if len(parts) == 3 {
		if parts[1] == "tag" {
			refType = vcs.Tag
		} else {
			refType = vcs.Branch
		}

		refName = parts[len(parts)-1]
	}

	commitInfo := &vcs.CommitInfo{
		Author:    event.Commits[0].Author.GetName(),
		CreatedAt: event.Commits[0].GetTimestamp().Time,
		Ref:       refName,
		RefType:   refType,
		Hash:      event.GetAfter(),
		Repository: domain.Repository{
			Provider: "github",
			Name:     event.Repo.GetName(),
			Owner:    event.Repo.GetOwner().GetName(),
		},
	}

	return commitInfo, nil
}
