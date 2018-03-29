package v1

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// SourceWebhook handles webhook request. It queues a build job if it is required.
func (a *Api) SourceWebhook(c echo.Context) error {
	usr, err := a.userStorage.Get(c.Param("username"))
	if err != nil {
		return err
	}

	commit, err := a.sourcesvc.ParseWebhookPayload(c.Param("provider"), c.Request())
	if err == domain.ErrSourceIrrelevantEvent {
		return c.NoContent(http.StatusOK)
	} else if err != nil {
		return err
	}

	imgKey := domain.ImageKey(commit.Repository)
	img, err := a.imageStorage.Get(usr.Username, imgKey)
	if err != nil {
		return err
	}

	tag, ok := img.MatchingTag(commit)
	if !ok {
		return c.NoContent(http.StatusOK)
	}

	tokens, err := a.oauthStorage.GetTokens(usr.Username)
	if err != nil {
		return err
	}

	token, ok := tokens[c.Param("provider")]
	if !ok {
		return domain.ErrAuthUnauthorized
	}

	job := domain.BuildJob{
		ImageKey:    imgKey,
		ImageName:   img.Name,
		Dockerfile:  img.DockerfilePath,
		Tag:         tag.Tag(commit),
		ImageRepo:   img.Repository,
		CommitHash:  commit.Hash,
		CommitRef:   commit.Ref,
		ImageOwner:  usr.Username,
		VcsToken:    token.Token,
		VcsUsername: token.Identity,
	}

	return a.buildsvc.Queue(job)
}

// SourceOrganisations responses with source client user's list of organisations
func (a *Api) SourceOrganisations(secrets domain.AuthSecrets, c echo.Context) error {
	provider := c.Param("provider")
	orgs, err := a.sourcesvc.Organisations(context.Background(), provider, secrets.Username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, orgs)
}

// SourceRepositories responses with source client user's list of repositories for
// given organisation
func (a *Api) SourceRepositories(secrets domain.AuthSecrets, c echo.Context) error {
	provider := c.Param("provider")
	org := c.QueryParam("org")

	repos, err := a.sourcesvc.Repositories(context.Background(), provider, secrets.Username, org)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, repos)
}
