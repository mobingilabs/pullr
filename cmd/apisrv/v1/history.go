package v1

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/storage"
	"github.com/pkg/errors"
)

func (a *API) imagesStatuses(username string, c echo.Context) error {
	listOpts := storage.NewListOptions()
	if err := c.Bind(listOpts); err != nil {
		listOpts = nil
	}

	statuses, err := a.statuses(username, "images", listOpts)
	if err != nil {
		return err
	}

	imgKeys := make([]string, len(statuses))
	for i, status := range statuses {
		imgKeys[i] = status.ID
	}

	images, err := a.Storage.GetImages(imgKeys)
	if err != nil {
		return errors.WithMessage(err, "failed to get images for statuses")
	}

	for i := range images {
		for j := range statuses {
			if images[i].Key == statuses[j].ID {
				images[i].Status = &statuses[j]
			}
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"images": images})
}

func (a *API) imageHistory(username string, c echo.Context) error {
	return a.history(username, "image", c.Param("key"), c)
}

func (a *API) statuses(username string, kind string, listOpts *storage.ListOptions) ([]domain.Status, error) {

	statuses, err := a.Storage.Statuses(username, kind, listOpts)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get statuses for %s", kind)
	}

	statuses = append([]domain.Status{}, statuses...)
	return statuses, nil
}

func (a *API) history(username string, kind string, id string, c echo.Context) error {
	statuses, err := a.Storage.History(username, kind, id, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to get statuses of %s#%s", kind, id)
	}

	statuses = append([]domain.Status{}, statuses...)
	return c.JSON(http.StatusOK, map[string]interface{}{"statuses": statuses})
}
