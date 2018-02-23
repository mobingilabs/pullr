package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo"
	"github.com/mobingilabs/pullr/pkg/domain"
	"github.com/mobingilabs/pullr/pkg/srv"
	"github.com/mobingilabs/pullr/pkg/storage"
	"github.com/mobingilabs/pullr/pkg/vcs"
	"github.com/mobingilabs/pullr/pkg/vcs/github"
	log "github.com/sirupsen/logrus"
)

func (a *API) gitWebhook(c echo.Context) error {
	provider := vcsFor(c.Param("provider"))
	if provider == nil {
		return c.NoContent(http.StatusNotFound)
	}

	webhook, err := provider.ParseWebhookRequest(c.Request())
	if err != nil {
		return err
	}

	if webhook.Event == vcs.EventPing {
		return c.NoContent(http.StatusOK)
	}

	if webhook.Event != vcs.EventPush {
		log.Warningf("Pullr doesn't support webhook events other than 'push', got '%s'.", webhook.Event)
		return srv.NewErrBadValue("event", webhook.Event)
	}

	commitInfo, err := provider.ExtractCommitInfo(webhook)
	if err != nil {
		return err
	}

	imgKey := domain.ImageKey(commitInfo.Repository)
	log.Infof("Got webhook event for image key '%s'", imgKey)

	img, err := a.Storage.FindImageByKey(imgKey)
	if err != nil {
		if err == storage.ErrNotFound {
			log.Warningf("Image not found with key '%s'", imgKey)
		}

		return err
	}

	dockerTag := getDockerTag(commitInfo, img.Tags)
	if dockerTag == "" {
		log.Infof("Push event doesn't match any image tags to build skipping...")
		return c.NoContent(http.StatusOK)
	}

	job := domain.NewBuildImageJob("pullr:apisrv", imgKey, commitInfo.Ref, commitInfo.Hash, dockerTag)
	jobData, err := json.Marshal(job)
	if err != nil {
		return err
	}

	log.Infof("Putting build image job to queue with image key '%s' and data '%s'", imgKey, string(jobData))
	if err := a.JobQ.Put(a.Conf.JobQ.BuildQueue, bytes.NewBuffer(jobData)); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func vcsFor(provider string) vcs.Vcs {
	switch provider {
	case "github":
		return github.New()
	default:
		return nil
	}
}

func getDockerTag(commit *vcs.CommitInfo, tags []domain.ImageTag) string {
	for _, t := range tags {
		// commit is a tag push
		if t.RefType == string(vcs.Tag) && commit.RefType == vcs.Tag {
			// Check for regexp tests
			if strings.HasPrefix(t.RefTest, "/") && strings.HasSuffix(t.RefTest, "/") {
				if len(t.RefTest) <= 2 {
					continue
				}

				rx, err := regexp.Compile(t.RefTest[1 : len(t.RefTest)-2])
				if err != nil || !rx.MatchString(commit.Ref) {
					continue
				}
			} else if t.RefTest != commit.Ref {
				continue
			}

			if t.Name != "" {
				return t.Name
			}

			return commit.Ref
		}

		// commit is a normal push on a branch
		if commit.Ref == t.RefTest {
			return t.Name
		}
	}

	// commit doesn't match any image tags
	return ""
}
