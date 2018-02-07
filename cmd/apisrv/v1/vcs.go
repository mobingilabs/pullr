package v1

import (
	"context"
	"net/http"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/vcs"
	"github.com/mobingilabs/pullr/pkg/vcs/github"
)

func (a *apiv1) VcsOrganisations(username string, c echo.Context) error {
	provider := c.Param("provider")
	usr, err := a.Storage.FindUser(username)
	if err != nil {
		return err
	}

	vcsClient := getVcsClient(provider, usr.Tokens[provider])
	if vcsClient == nil {
		glog.Errorf("Unsupported vcs provider '%s'", c.Param("provider"))
		return c.NoContent(http.StatusBadRequest)
	}

	orgs, err := vcsClient.ListOrganisations(context.Background())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string][]string{"organisations": orgs})
}

func (a *apiv1) VcsRepositories(username string, c echo.Context) error {
	provider := c.Param("provider")
	usr, err := a.Storage.FindUser(username)
	if err != nil {
		return err
	}

	vcsClient := getVcsClient(provider, usr.Tokens[provider])
	if vcsClient == nil {
		glog.Errorf("Unsupported vcs provider '%s'", c.Param("provider"))
	}

	organisation := c.Param("organisation")
	repos, err := vcsClient.ListRepositories(context.Background(), organisation)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string][]string{"repositories": repos})
}

func getVcsClient(provider, token string) vcs.Client {
	switch provider {
	case "github":
		return github.NewClientWithToken(token)
	default:
		return nil
	}
}
