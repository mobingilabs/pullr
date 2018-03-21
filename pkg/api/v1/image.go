package v1

import (
	"net/http"
	"strings"
	"time"

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

	// Ignore owner & key fields
	img.Key = domain.ImageKey(img.Repository)
	img.Owner = secrets.Username
	img.CreatedAt = time.Now()
	img.UpdatedAt = img.CreatedAt

	valid, err := img.Valid()
	if !valid {
		return err
	}

	_, err = a.imageStorage.Get(secrets.Username, img.Key)
	if err == nil {
		return domain.ErrImageExists
	}

	err = a.imageStorage.Put(secrets.Username, img)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, img)
}

// ImageGet responds with the image details found by the :key parameter.
func (a *Api) ImageGet(secrets domain.AuthSecrets, c echo.Context) error {
	imgKey := strings.TrimSpace(c.Param("key"))
	if imgKey == "" {
		return domain.ErrNotFound
	}

	img, err := a.imageStorage.Get(secrets.Username, imgKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, img)
}

// ImageUpdate accepts partial domain.Image as it is body and updates the
// matching image record found in storage with the body
func (a *Api) ImageUpdate(secrets domain.AuthSecrets, c echo.Context) error {
	imgKey := strings.TrimSpace(c.Param("key"))
	if imgKey == "" {
		return domain.ErrNotFound
	}

	var update domain.Image
	if err := c.Bind(&update); err != nil {
		return err
	}

	orig, err := a.imageStorage.Get(secrets.Username, imgKey)
	if err != nil {
		return err
	}

	update.Key = domain.ImageKey(update.Repository)
	update.Owner = secrets.Username
	update.CreatedAt = orig.CreatedAt
	update.UpdatedAt = time.Now()

	valid, err := update.Valid()
	if !valid {
		return err
	}

	err = a.imageStorage.Update(secrets.Username, imgKey, update)
	if err != nil {
		return err
	}

	update, err = a.imageStorage.Get(secrets.Username, update.Key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, update)
}

// ImageDelete deletes the image found by it's key found in the url.
func (a *Api) ImageDelete(secrets domain.AuthSecrets, c echo.Context) error {
	imgKey := strings.TrimSpace(c.Param("key"))
	if imgKey == "" {
		return domain.ErrNotFound
	}

	return a.imageStorage.Delete(secrets.Username, imgKey)
}
