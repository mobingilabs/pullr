package github

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/vcs"
)

type githubVcs struct{}

// New creates github specific VCS implementation
func New() vcs.Vcs {
	return &githubVcs{}
}

func (*githubVcs) ExtractCommitInfo(r *vcs.WebhookRequest) (*vcs.CommitInfo, error) {
	switch r.Event {
	case vcs.PushEvent:
		return extractCommitInfoPushPayload(r)
	default:
		return nil, vcs.ErrUnsupportedEvent
	}
}

func (*githubVcs) ParseWebhookRequest(r *http.Request) (*vcs.WebhookRequest, error) {
	if !strings.HasPrefix(r.Header.Get("User-Agent"), "GitHub-Hookshot") {
		return nil, vcs.ErrInvalidWebhook
	}

	event := r.Header.Get("X-githubVcs-Event")
	if event == "" {
		return nil, vcs.ErrInvalidWebhook
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer errs.Log(r.Body.Close())

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
	refParts := strings.Split(event.GetRef(), "/")
	if len(refParts) != 3 {
		return nil, vcs.ErrInvalidWebhookPayload
	}

	var refType vcs.RefType
	if refParts[1] == "tag" {
		refType = vcs.Tag
	} else {
		refType = vcs.Branch
	}

	refName := refParts[len(refParts)-1]

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
