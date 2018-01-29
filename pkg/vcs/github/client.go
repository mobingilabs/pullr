package github

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/mobingilabs/pullr/pkg/domain"
)

// Client version control system
type Client struct {
	Token *string
}

// NewClient creates an unauthenticated Client
func NewClient() *Client {
	return &Client{}
}

// NewClientWithToken creates an authenticated Client
func NewClientWithToken(token string) *Client {
	return &Client{Token: &token}
}

// CheckFileExists checks if the given path is exists on the repository
func (g *Client) CheckFileExists(ctx context.Context, repository *domain.Repository, path string, ref string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", repository.Owner, repository.Name, path, ref)

	request, err := http.NewRequest("HEAD", url, nil)
	request = request.WithContext(ctx)
	if err != nil {
		return false, err
	}

	if g.Token != nil {
		request.Header.Add("Authorization", fmt.Sprintf("token %s", *g.Token))
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
func (g *Client) CloneRepository(ctx context.Context, repository *domain.Repository, clonePath string, ref string) (string, error) {
	// FIXME: This will save token to .git/config, possible security risk, altough we gonna remove the directory after build
	var cloneURL string
	if g.Token != nil {
		cloneURL = fmt.Sprintf("https://%s@github.com/%s/%s.git", *g.Token, repository.Owner, repository.Name)
	} else {
		cloneURL = fmt.Sprintf("https://git@github.com/%s/%s.git", repository.Owner, repository.Name)
	}

	cloneCmd := exec.Command("git", "clone", cloneURL, clonePath)
	if err := cloneCmd.Run(); err != nil {
		return "", err
	}

	checkoutCmd := exec.CommandContext(ctx, "git", "checkout", ref)
	checkoutCmd.Dir = clonePath

	stdoutReader, err := checkoutCmd.StdoutPipe()
	if err != nil {
		return "", checkoutCmd.Run()
	}

	cmdErr := checkoutCmd.Run()

	var logBytes []byte
	if _, err = stdoutReader.Read(logBytes); err != nil {
		logBytes = []byte{' '}
	}

	return string(logBytes), cmdErr
}
