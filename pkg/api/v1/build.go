package v1

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
)

// BuildList response with list of images sorted by their last build times.
// ListOptions can be used for paginating the results.
func (a *Api) BuildList(secrets domain.AuthSecrets, c echo.Context) error {
	type responsePayload struct {
		Builds     []domain.Build          `json:"builds"`
		Images     map[string]domain.Image `json:"images"`
		Pagination domain.Pagination       `json:"pagination"`
	}

	response := responsePayload{
		Builds: []domain.Build{},
		Images: make(map[string]domain.Image),
	}

	listOpts := domain.DefaultListOptions
	_ = c.Bind(&listOpts)

	builds, pagination, err := a.buildStorage.List(secrets.Username, listOpts)
	if err == domain.ErrNotFound {
		return c.JSON(http.StatusOK, response)
	} else if err != nil {
		return err
	}

	imgKeys := make([]string, len(builds))
	for i := range builds {
		imgKeys[i] = builds[i].ImageKey
	}

	imgs, err := a.imageStorage.GetMany(secrets.Username, imgKeys)
	if err != nil {
		return err
	}

	response.Builds = builds
	response.Images = imgs
	response.Pagination = pagination
	return c.JSON(http.StatusOK, response)
}

// BuildHistory response with history of builds of an image. ListOptions
// can be used for paginating the results.
func (a *Api) BuildHistory(secrets domain.AuthSecrets, c echo.Context) error {
	type responsePayload struct {
		BuildRecords []domain.BuildRecord `json:"build_records"`
		Pagination   domain.Pagination    `json:"pagination"`
	}

	listOpts := domain.DefaultListOptions
	_ = c.Bind(&listOpts)

	imgKey := strings.TrimSpace(c.Param("key"))
	if imgKey == "" {
		return domain.ErrNotFound
	}

	records, pagination, err := a.buildStorage.GetAll(secrets.Username, imgKey, listOpts)
	if err == domain.ErrNotFound {
		return c.JSON(http.StatusOK, responsePayload{[]domain.BuildRecord{}, domain.Pagination{}})
	} else if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, responsePayload{records, pagination})
}
