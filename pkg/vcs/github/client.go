package github

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/vcs"
	"golang.org/x/oauth2"
	"gopkg.in/src-d/go-git.v4/plumbing"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// Client encapsulates authenticated GithubAPI requests
type Client struct {
	Username string
	Token    string
}

// NewClientWithToken creates an authenticated GithubAPI client
func NewClientWithToken(username, token string) *Client {
	return &Client{Username: username, Token: token}
}

// CheckFileExists checks a repository if the given file path exists or not
func (g *Client) CheckFileExists(ctx context.Context, repository domain.Repository, path string, ref string) (bool, error) {
	cl := newAuthenticatedClient(ctx, g.Token)
	reader, err := cl.Repositories.DownloadContents(ctx, repository.Owner, repository.Name, path, &github.RepositoryContentGetOptions{Ref: ref})
	if err != nil {
		return false, nil
	}
	errs.Log(reader.Close())
	return true, nil
}

// CloneRepository clones the repository content to given path
func (g *Client) CloneRepository(ctx context.Context, repository domain.Repository, clonePath string, ref string) error {
	repo, err := git.PlainClone(clonePath, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: g.Username,
			Password: g.Token,
		},
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

// ListOrganisations fetches authenticated user's organisations
func (g *Client) ListOrganisations(ctx context.Context) ([]string, error) {
	if g.Token == "" {
		return nil, vcs.ErrAuthRequired
	}

	client := newAuthenticatedClient(ctx, g.Token)
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

// ListRepositories fetches authenticated user's repositories
func (g *Client) ListRepositories(ctx context.Context, organisation string) ([]string, error) {
	if g.Token == "" {
		return nil, vcs.ErrAuthRequired
	}

	client := newAuthenticatedClient(ctx, g.Token)
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

func newAuthenticatedClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
