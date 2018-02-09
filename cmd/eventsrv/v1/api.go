package v1

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
	"github.com/mobingilabs/pullr/pkg/comm"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/storage"
	"github.com/mobingilabs/pullr/pkg/vcs/github"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/vcs"
)

type apiv1 struct {
	Storage storage.Storage
	Queue   comm.JobTransporter
	Group   *echo.Group
}

func (a *apiv1) webhookHandler(c echo.Context) error {
	provider := vcsFor(c.Param("provider"))
	if provider == nil {
		return c.NoContent(http.StatusNotFound)
	}

	webhook, err := provider.ParseWebhookRequest(c.Request())
	if err != nil {
		return err
	}

	if webhook.Event != vcs.PushEvent {
		glog.Warningf("Pullr doesn't support webhook events other than 'push', got '%s'.", webhook.Event)
		return c.NoContent(http.StatusNoContent)
	}

	commitInfo, err := provider.ExtractCommitInfo(webhook)
	if err != nil {
		return err
	}

	imgKey := domain.ImageKey(commitInfo.Repository)
	glog.Infof("Got webhook event for image key '%s'", imgKey)

	_, err = a.Storage.FindImageByKey(imgKey)
	if err != nil {
		if err == storage.ErrNotFound {
			glog.Warningf("Image not found with key '%s'", imgKey)
		}

		return err
	}

	job := domain.NewBuildImageJob("eventsrv", imgKey)
	jobData, err := json.Marshal(job)
	if err != nil {
		return err
	}

	glog.Infof("Putting build image job to queue with image key '%s' and data '%s'", imgKey, string(jobData))
	return a.Queue.Put(domain.BuildQueue, bytes.NewBuffer(jobData))
}

func vcsFor(provider string) vcs.Vcs {
	switch provider {
	case "github":
		return github.New()
	default:
		return nil
	}
}

func NewApiV1(e *echo.Echo, storage storage.Storage, queue comm.JobTransporter) *apiv1 {
	g := e.Group("/v1")

	// TODO: Handle other storage options
	api := &apiv1{
		Group:   g,
		Storage: storage,
		Queue:   queue,
	}

	g.POST("/:provider", api.webhookHandler)

	return api
}
