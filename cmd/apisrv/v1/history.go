package v1

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/storage"
)

func (a *API) statuses(username string, c echo.Context) error {
	resourceKind := c.Param("kind")

	listOpts := new(storage.ListOptions)
	if err := c.Bind(listOpts); err != nil {
		listOpts = nil
	}

	statuses, err := a.Storage.Statuses(username, resourceKind, listOpts)
	if err != nil {
		return err
	}

	statuses = append([]domain.Status{}, statuses...)
	return c.JSON(http.StatusOK, map[string]interface{}{"statuses": statuses})
}
