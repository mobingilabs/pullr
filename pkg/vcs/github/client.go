package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/vcs"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type githubClient struct {
	username string
	token    string
}

// NewClientWithToken creates an authenticated github specific vcs client
func NewClientWithToken(username, token string) vcs.Client {
	return &githubClient{username: username, token: token}
}

func (g *githubClient) CheckFileExists(ctx context.Context, repository domain.Repository, path string, ref string) (bool, error) {
	cl := newAuthenticatedClient(ctx, g.token)
	reader, err := cl.Repositories.DownloadContents(ctx, repository.Owner, repository.Name, path, &github.RepositoryContentGetOptions{Ref: ref})
	if err != nil {
		return false, nil
	}
	errs.Log(reader.Close())
	return true, nil
}

func (g *githubClient) CloneRepository(ctx context.Context, repository domain.Repository, clonePath string, ref string) error {
	repo, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: g.username,
			Password: g.token,
		},
		URL: fmt.Sprintf("https://github.com/%s/%s", repository.Owner, repository.Name),
	})
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	return wt.Checkout(&git.CheckoutOptions{Hash: plumbing.NewHash(ref)})
}

func (g *githubClient) ListOrganisations(ctx context.Context) ([]string, error) {
	if g.token == "" {
		return nil, vcs.ErrAuthRequired
	}

	client := newAuthenticatedClient(ctx, g.token)
	usr, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	orgs, _, err := client.Organizations.List(ctx, "", &github.ListOptions{
		PerPage: 100,
	})
	if err != nil {
		return nil, err
	}

	orgNames := make([]string, len(orgs)+1)
	orgNames[0] = usr.GetLogin()
	for i, org := range orgs {
		if org.GetLogin() == "" {
			continue
		}

		orgNames[i+1] = org.GetLogin()
	}

	return orgNames, nil
}

func (g *githubClient) ListRepositories(ctx context.Context, organisation string) ([]string, error) {
	if g.token == "" {
		return nil, vcs.ErrAuthRequired
	}

	client := newAuthenticatedClient(ctx, g.token)
	usr, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}

	// Otherwise try to get organisation repositories
	var repos []string
	page := 0
	lastPage := 100
	for page <= lastPage {
		var gRepos []*github.Repository
		var res *github.Response
		var err error

		if organisation == usr.GetLogin() {
			gRepos, res, err = client.Repositories.List(ctx, "", &github.RepositoryListOptions{
				Sort:        "name",
				Affiliation: "owner",
				ListOptions: github.ListOptions{
					PerPage: 100,
					Page:    page,
				},
			})
		} else {
			gRepos, res, err = client.Repositories.ListByOrg(ctx, organisation, &github.RepositoryListByOrgOptions{
				Type: "member",
				ListOptions: github.ListOptions{
					PerPage: 100,
					Page:    page,
				},
			})
		}
		if err != nil {
			return nil, err
		}

		pageRepos := make([]string, len(gRepos))
		for i, r := range gRepos {
			pageRepos[i] = r.GetName()
		}

		repos = append(repos, pageRepos...)

		if page == res.NextPage {
			break
		}

		page = res.NextPage
		lastPage = res.LastPage
	}

	return repos, nil
}

func (g *githubClient) RegisterWebhook(ctx context.Context, repo domain.Repository, uri string) error {
	c := newAuthenticatedClient(ctx, g.token)
	_, _, err := c.Repositories.CreateHook(ctx, repo.Owner, repo.Name, &github.Hook{
		Name: github.String("web"),
		Config: map[string]interface{}{
			"url":          github.String(uri),
			"content_type": "json",
		},
		Active: github.Bool(true),
		Events: []string{"push"},
	})
	return err
}

func newAuthenticatedClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
