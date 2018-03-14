package v1

import (
	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// ImageList response with list of user's images. ListOptions can be used
// for pagination
func (a *Api) ImageList(secrets domain.AuthSecrets, c echo.Context) error {
	return nil
}

// ImageCreates accepts domain.Image as it is body and stores accepted image
// data in images storage.
func (a *Api) ImageCreate(secrets domain.AuthSecrets, c echo.Context) error {
	return nil
}

// ImageGet responds with the image details found by the :key parameter.
func (a *Api) ImageGet(secrets domain.AuthSecrets, c echo.Context) error {
	return nil
}

// ImageUpdate accepts partial domain.Image as it is body and updates the
// matching image record found in storage with the body
func (a *Api) ImageUpdate(secret domain.AuthSecrets, c echo.Context) error {
	return nil
}

// ImageDelete deletes the image found by it's key found in the url.
func (a *Api) ImageDelete(secrets domain.AuthSecrets, c echo.Context) error {
	return nil
}
