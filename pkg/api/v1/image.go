package v1

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// ImageList response with list of user's images. ListOptions can be used
// for pagination
func (a *Api) ImageList(secrets domain.AuthSecrets, c echo.Context) error {
	type responsePayload struct {
		Images     []domain.Image    `json:"images"`
		Pagination domain.Pagination `json:"pagination"`
	}

	listOpts := domain.DefaultListOptions
	_ = c.Bind(&listOpts)

	imgs, pagination, err := a.imageStorage.List(secrets.Username, listOpts)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, responsePayload{imgs, pagination})
}

// ImageCreates accepts domain.Image as it is body and stores accepted image
// data in images storage.
func (a *Api) ImageCreate(secrets domain.AuthSecrets, c echo.Context) error {
	var img domain.Image
	if err := c.Bind(&img); err != nil {
		return err
	}

	// Ignore the owner from the body
	img.Owner = secrets.Username

	valid, validationErrs := img.Valid()
	if !valid {
		return validationErrs
	}

	// Ignore the key field from the body
	img.Key = domain.ImageKey(img)

	err := a.imageStorage.Put(secrets.Username, img)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, img)
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
