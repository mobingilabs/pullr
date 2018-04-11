package github

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/mobingilabs/pullr/pkg/domain"
)

const apiURL = "https://api.github.com"

// Client implements domain.SourceClient. Can parse webhooks, and
// query authenticated github user's repositories
type Client struct{}

// NewClient creates a github client
func NewClient() *Client {
	return &Client{}
}

func (*Client) doRequest(ctx context.Context, apiReq apiRequest) (int, []byte, error) {
	if apiReq.params == nil {
		apiReq.params = url.Values{}
	}
	apiReq.params.Set("access_token", apiReq.accessToken)
	if apiReq.method == "" {
		apiReq.method = http.MethodGet
	}

	reqURL := fmt.Sprintf("%s%s?%s", apiURL, apiReq.path, apiReq.params.Encode())
	req, err := http.NewRequest(apiReq.method, reqURL, apiReq.body)
	if err != nil {
		return 0, nil, err
	}

	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, nil, err
	}

	return res.StatusCode, body, nil
}

// RegisterWebhook registers pullr to source provider's webhooks
func (c *Client) RegisterWebhook(ctx context.Context, token string, webhookURL string, repo domain.SourceRepository) error {
	type registerConfig struct {
		Url         string `json:"url"`
		ContentType string `json:"content_type"`
	}
	type registerBody struct {
		Name   string         `json:"name"`
		Config registerConfig `json:"config"`
	}

	body := registerBody{
		Name: "web",
		Config: registerConfig{
			Url:         webhookURL,
			ContentType: "json",
		},
	}
	var bodyJson bytes.Buffer
	err := json.NewEncoder(&bodyJson).Encode(body)
	if err != nil {
		return err
	}

	req := apiRequest{
		body:        &bodyJson,
		accessToken: token,
		method:      http.MethodPost,
		path:        fmt.Sprintf("/repos/%s/%s/hooks", repo.Owner, repo.Name),
	}
	code, resBody, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}
	if code != http.StatusCreated {
		return errors.New(string(resBody))
	}

	return nil
}

// ParseWebhookPayload parses github's webhook request and extracts commit info out of it
func (*Client) ParseWebhookPayload(req *http.Request) (*domain.CommitInfo, error) {
	if !strings.HasPrefix(req.Header.Get("User-Agent"), "GitHub-Hookshot") {
		return nil, domain.ErrSourceBadPayload
	}

	event := req.Header.Get("X-GitHub-Event")
	if event == "" {
		return nil, domain.ErrSourceBadPayload
	}

	if event != "push" {
		return nil, domain.ErrSourceIrrelevantEvent
	}

	body, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		return nil, err
	}

	var pushEvent PushEvent
	err = json.Unmarshal(body, &pushEvent)
	if err != nil {
		return nil, err
	}

	valid, err := pushEvent.Validate()
	if !valid {
		return nil, err
	}

	// Parse push type
	refParts := strings.Split(*pushEvent.Ref, "/")
	if len(refParts) != 3 {
		return nil, domain.ErrSourceBadPayload
	}

	var refType domain.SourceRefType
	if refParts[1] == "tags" {
		refType = domain.SourceTag
	} else {
		refType = domain.SourceBranch
	}

	refName := refParts[len(refParts)-1]
	commit := pushEvent.HeadCommit

	commitInfo := &domain.CommitInfo{
		Author:    *commit.Author.Name,
		CreatedAt: *commit.Timestamp,
		Ref:       refName,
		RefType:   refType,
		Hash:      *pushEvent.After,
		Repository: domain.SourceRepository{
			Provider: "github",
			Name:     *pushEvent.Repository.Name,
			Owner:    *pushEvent.Repository.Owner.Login,
		},
	}

	return commitInfo, nil
}

// Organisations reports back user's organisations
func (c *Client) Organisations(ctx context.Context, identity string, token string) ([]string, error) {
	req := apiRequest{
		path:        "/user/orgs",
		params:      url.Values{"state": {"active"}, "per_page": {"500"}},
		accessToken: token,
	}

	code, body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var orgs []struct {
		Login string `json:"login"`
	}
	if err := json.Unmarshal(body, &orgs); err != nil {
		return nil, err
	}

	orgNames := make([]string, 0, len(orgs)+1)
	orgNames = append(orgNames, identity)
	for _, org := range orgs {
		orgNames = append(orgNames, org.Login)
	}

	return orgNames, nil
}

// Repositories reports back given organisation's repositories. If you want to
// get user's own repositories pass username as organisation
func (c *Client) Repositories(ctx context.Context, identity string, organisation string, token string) ([]string, error) {
	path := "/user/repos"
	params := url.Values{"affiliation": {"owner,collaborator"}}
	if organisation != identity {
		path = fmt.Sprintf("/orgs/%s/repos", organisation)
		params = url.Values{"type": {"member"}}
	}

	// Avoid pagination
	params.Set("per_page", "500")

	req := apiRequest{
		path:        path,
		params:      params,
		accessToken: token,
	}

	code, body, err := c.doRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var repoList []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(body, &repoList); err != nil {
		return nil, err
	}

	repoNames := make([]string, 0, len(repoList))
	for _, repo := range repoList {
		repoNames = append(repoNames, repo.Name)
	}

	return repoNames, nil
}

func (c *Client) identity(ctx context.Context, token string) (string, error) {
	req := apiRequest{
		path:        "/user",
		accessToken: token,
	}

	code, body, err := c.doRequest(ctx, req)
	if err != nil {
		return "", err
	}
	if code != 200 {
		return "", errors.New(string(body))
	}

	var profile struct {
		Login string `json:"login"`
	}
	if err := json.Unmarshal(body, &profile); err != nil {
		return "", err
	}

	return profile.Login, nil
}

type apiRequest struct {
	method      string
	path        string
	accessToken string
	params      url.Values
	body        io.Reader
}
