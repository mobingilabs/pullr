package v1

import (
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// SourceWebhook handles webhook request accept source control provider
// webhook as it is request body. It queues a build job if it is required.
func (a *Api) SourceWebhook(c echo.Context) error {
	return nil
}

// SourceOrganisations responses with source client user's list of organisations
func (a *Api) SourceOrganisations(secrets domain.AuthSecrets, c echo.Context) error {
	return nil
}

// SourceRepositories responses with source client user's list of repositories for
// given organisation
func (a *Api) SourceRepositories(secrets domain.AuthSecrets, c echo.Context) error {
	return nil
}
