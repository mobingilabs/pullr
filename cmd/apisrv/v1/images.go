package v1

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

func (a *apiv1) imagesIndex(user string, c echo.Context) error {
	images, err := a.Storage.FindAllImages(user)
	if err != nil {
		return err
	}

	if images == nil {
		return c.JSON(http.StatusOK, []domain.Image{})
	}

	return c.JSON(http.StatusOK, images)
}

func (a *apiv1) imagesCreate(user string, c echo.Context) error {
	type createPayload struct {
		domain.Image
		Tags []domain.ImageTag `json:"tags"`
	}

	payload := new(createPayload)
	if err := c.Bind(payload); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	payload.Image.Owner = user
	if err := a.Storage.CreateImage(payload.Image); err != nil {
		return err
	}

	imageKey := domain.ImageKey(payload.Image.Repository)
	for _, tag := range payload.Tags {
		if err := a.Storage.CreateImageTag(imageKey, tag); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusNoContent)
}

func (a *apiv1) imagesDelete(user string, c echo.Context) error {
	imageKey := strings.TrimSpace(c.Param("key"))
	if imageKey == "" {
		return c.NoContent(http.StatusNotFound)
	}

	image, err := a.Storage.FindImageByKey(imageKey)
	if err != nil {
		return err
	}

	if image.Owner != user {
		return c.NoContent(http.StatusNotFound)
	}

	if err := a.Storage.DeleteImage(imageKey); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (a *apiv1) imagesUpdate(user string, c echo.Context) error {
	type updatePayload struct {
		domain.Image
		Tags []domain.ImageTag `json:"tags"`
	}

	key := strings.TrimSpace(c.Param("key"))
	if key == "" {
		return c.NoContent(http.StatusNotFound)
	}

	payload := new(updatePayload)
	if err := c.Bind(payload); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	a.Storage.UpdateImage(key, payload.Image)

	newImageKey := domain.ImageKey(payload.Image.Repository)
	return c.JSON(http.StatusOK, struct {
		ImageKey string `json:"image_key"`
	}{newImageKey})
}
