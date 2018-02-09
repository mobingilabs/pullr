package github

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/google/go-github/github"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/vcs"
	"golang.org/x/oauth2"
)

// Client version control system
type Client struct {
	Token string
}

// NewClient creates an unauthenticated Client
func NewClient() *Client {
	return &Client{}
}

// NewClientWithToken creates an authenticated Client
func NewClientWithToken(token string) *Client {
	return &Client{Token: token}
}

// CheckFileExists checks if the given path is exists on the repository
func (g *Client) CheckFileExists(ctx context.Context, repository domain.Repository, path string, ref string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", repository.Owner, repository.Name, path, ref)

	request, err := http.NewRequest("HEAD", url, nil)
	request = request.WithContext(ctx)
	if err != nil {
		return false, err
	}

	if g.Token != "" {
		request.Header.Add("Authorization", fmt.Sprintf("token %s", g.Token))
	}

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	// http status code 200 means the file is exists
	return res.StatusCode == 200, nil
}

// CloneRepository will clone repository on the disk in given path
func (g *Client) CloneRepository(ctx context.Context, repository domain.Repository, clonePath string, ref string) ([]byte, error) {
	// FIXME: This will save token to .git/config, possible security risk, altough we gonna remove the directory after build
	var cloneURL string
	if g.Token != "" {
		cloneURL = fmt.Sprintf("https://%s@github.com/%s/%s.git", g.Token, repository.Owner, repository.Name)
	} else {
		cloneURL = fmt.Sprintf("https://git@github.com/%s/%s.git", repository.Owner, repository.Name)
	}

	cloneCmd := exec.Command("git", "clone", cloneURL, clonePath)
	if err := cloneCmd.Run(); err != nil {
		return nil, err
	}

	checkoutCmd := exec.CommandContext(ctx, "git", "checkout", ref)
	checkoutCmd.Dir = clonePath

	stdoutReader, err := checkoutCmd.StdoutPipe()
	if err != nil {
		return nil, checkoutCmd.Run()
	}

	cmdErr := checkoutCmd.Run()

	var logBytes []byte
	if _, err = stdoutReader.Read(logBytes); err != nil {
		logBytes = nil
	}

	return logBytes, cmdErr
}

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
