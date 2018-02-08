package v1

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/storage"
)

func (a *apiv1) imagesGet(user string, c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	image, err := a.Storage.FindImageByKey(key)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, image)
}

func (a *apiv1) imagesIndex(username string, c echo.Context) (err error) {
	var images []domain.Image
	var pagination storage.Pagination
	if sinceQuery := c.QueryParam("since"); sinceQuery != "" {
		i, err := strconv.ParseInt(sinceQuery, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid time format"})
		}

		since := time.Unix(i, 0)
		images, err = a.Storage.FindAllImagesSince(username, since)
	} else {
		listOpts := new(storage.ListOptions)
		if err := c.Bind(listOpts); err != nil {
			listOpts = nil
		}

		images, pagination, err = a.Storage.FindAllImages(username, listOpts)
	}

	if err != nil {
		return err
	}

	if images == nil {
		images = []domain.Image{}
	}

	res := struct {
		Images     []domain.Image     `json:"images"`
		Pagination storage.Pagination `json:"pagination"`
	}{images, pagination}

	return c.JSON(http.StatusOK, res)
}

func (a *apiv1) imagesCreate(user string, c echo.Context) error {
	payload := new(domain.Image)
	if err := c.Bind(payload); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	payload.Owner = user
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = payload.CreatedAt
	imageKey, err := a.Storage.CreateImage(*payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"key": imageKey})
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
	key := strings.TrimSpace(c.Param("key"))
	if key == "" {
		return c.NoContent(http.StatusNotFound)
	}

	payload := new(domain.Image)
	if err := c.Bind(payload); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	payload.UpdatedAt = time.Now()
	payload.CreatedAt = time.Time{} // don't allow payload to update creation time
	newKey, err := a.Storage.UpdateImage(key, *payload)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"key": newKey})
}
