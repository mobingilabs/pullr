package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/errs"
	"github.com/mobingilabs/pullr/pkg/srv"
	"github.com/mobingilabs/pullr/pkg/storage"
	"github.com/mobingilabs/pullr/pkg/vcs/github"
	"github.com/sirupsen/logrus"
)

func (a *API) imagesGet(user string, c echo.Context) error {
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

func (a *API) imagesIndex(username string, c echo.Context) error {
	type indexResponse struct {
		Images     []domain.Image     `json:"images"`
		Pagination storage.Pagination `json:"pagination"`
	}

	if sinceQuery := c.QueryParam("since"); sinceQuery != "" {
		i, err := strconv.ParseInt(sinceQuery, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid time format"})
		}

		since := time.Unix(i, 0)
		images, err := a.Storage.FindAllImagesSince(username, since)
		if err != nil {
			return err
		}

		// Make sure it is an empty array instead of nil
		images = append([]domain.Image{}, images...)

		return c.JSON(http.StatusOK, indexResponse{images, storage.Pagination{}})

	}

	listOpts := new(storage.ListOptions)
	if err := c.Bind(listOpts); err != nil {
		listOpts = nil
	}

	images, pagination, err := a.Storage.FindAllImages(username, listOpts)
	if err != nil {
		return err
	}

	// Make sure it is an empty array instead of nil
	images = append([]domain.Image{}, images...)

	return c.JSON(http.StatusOK, indexResponse{images, pagination})
}

func (a *API) imagesCreate(username string, c echo.Context) error {
	img := new(domain.Image)
	if err := c.Bind(img); err != nil {
		return srv.NewErrBadValue("body", "Invalid image structure")
	}

	if err := a.validateNewImg(img); err != nil {
		return err
	}

	img.Owner = username
	img.CreatedAt = time.Now()
	img.UpdatedAt = img.CreatedAt

	if strings.TrimSpace(img.DockerfilePath) == "" {
		img.DockerfilePath = "./Dockerfile"
	}

	// Check oauth token for the img repository before saving it
	usr, err := a.Storage.FindUser(username)
	if err != nil {
		return err
	}

	oauthToken := usr.Token(img.Repository.Provider)
	if oauthToken == nil {
		return srv.NewErrBadRequest(map[string]interface{}{
			"repository.provider": "Unauthenticated vcs provider",
		})
	}

	imgKey, err := a.Storage.CreateImage(*img)
	if err != nil {
		return err
	}

	// Register pullr webhook on vcs provider
	whURL := fmt.Sprintf("https://%s/api/v1/webhook/%s", c.Request().Host, img.Repository.Provider)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*30)
	defer cancel()

	logrus.WithFields(logrus.Fields{
		"user": username,
		"repo": img.Repository,
		"url":  whURL,
	}).Infof("registering webhook")

	g := github.NewClientWithToken(oauthToken.Username, oauthToken.Token)
	if err := g.RegisterWebhook(ctx, img.Repository, whURL); err != nil {
		// TODO: should we proceed without webhook?
		// Remove the image record if we failed to create webhook

		errs.Log(a.Storage.DeleteImage(imgKey))

		logrus.WithFields(logrus.Fields{
			"user": username,
			"repo": img.Repository,
		}).Errorf("Failed to register webhook: %v", err)
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"key": imgKey})
}

func (a *API) imagesDelete(user string, c echo.Context) error {
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

func (a *API) imagesUpdate(user string, c echo.Context) error {
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

func (a *API) validateNewImg(img *domain.Image) error {
	mistakes := make(map[string]interface{})

	if strings.TrimSpace(img.Name) == "" {
		mistakes["name"] = "Can't be empty"
	}

	if _, ok := a.OAuth[img.Repository.Provider]; !ok {
		mistakes["repository.provider"] = "Unsupported vcs provider"
	}

	if len(img.Tags) == 0 {
		mistakes["tags"] = "At least one docker tag needs to exist"
	}

	if len(mistakes) > 0 {
		return srv.NewErrBadRequest(mistakes)
	}

	return nil
}
