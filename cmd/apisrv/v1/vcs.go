package v1

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/srv"
	"github.com/mobingilabs/pullr/pkg/vcs"
	"github.com/mobingilabs/pullr/pkg/vcs/github"
)

func (a *API) vcsOrganisations(username string, c echo.Context) error {
	provider := c.Param("provider")
	usr, err := a.Storage.FindUser(username)
	if err != nil {
		return err
	}

	token := usr.Token(provider)
	if token == nil {
		return srv.NewErr("ERR_OAUTH_LOGIN", http.StatusUnauthorized, "OAuth token for requested vcs provider does not exist")
	}

	vcsClient := getVcsClient(provider, token)
	if vcsClient == nil {
		return srv.NewErrUnsupported("VCS provider '%s'", provider)
	}

	orgs, err := vcsClient.ListOrganisations(context.Background())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string][]string{"organisations": orgs})
}

func (a *API) vcsRepositories(username string, c echo.Context) error {
	provider := c.Param("provider")
	usr, err := a.Storage.FindUser(username)
	if err != nil {
		return err
	}

	token := usr.Token(provider)
	if token == nil {
		return srv.NewErr("ERR_OAUTH_LOGIN", http.StatusUnauthorized, "OAuth token for requested vcs provider does not exist")
	}

	vcsClient := getVcsClient(provider, token)
	if vcsClient == nil {
		return srv.NewErrUnsupported("VCS provider '%s'", provider)
	}

	organisation := c.Param("organisation")
	repos, err := vcsClient.ListRepositories(context.Background(), organisation)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string][]string{"repositories": repos})
}

func getVcsClient(provider string, token *domain.UserToken) vcs.Client {
	switch provider {
	case "github":
		return github.NewClientWithToken(token.Username, token.Token)
	default:
		return nil
	}
}
