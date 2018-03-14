package v1

import (
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// BuildList response with list of images sorted by their last build times.
// ListOptions can be used for paginating the results.
func (a *Api) BuildList(secrets domain.AuthSecrets, c echo.Context) error {
	return nil
}

// BuildHistory response with history of builds of an image. ListOptions
// can be used for paginating the results.
func (a *Api) BuildHistory(secrets domain.AuthSecrets, c echo.Context) error {
	return nil
}
