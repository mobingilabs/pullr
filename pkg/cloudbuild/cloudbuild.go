package cloudbuild

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/mobingilabs/pullr/pkg/domain"
)

const buildScript = `
docker build -f $PULLR_DOCKERFILE -t $PULLR_REGISTRY/$PULLR_OWNER/$PULLR_NAME:$PULLR_TAG .;
docker login -u $PULLR_REGISTRY_USER -p $PULLR_REGISTRY_PASSWORD $PULLR_REGISTRY;
docker push $PULLR_REGISTRY/$PULLR_OWNER/$PULLR_NAME:$PULLR_TAG;
`

const buildSpecTemplate = `
version: 0.2

phases:
  build:
    commands:
      - sh -c "%s"
`

// Pipeline builds and pushes docker images to registries by using aws cloudbuild service
type Pipeline struct {
	cb       *codebuild.CodeBuild
	logs     *cloudwatchlogs.CloudWatchLogs
	registry string
}

func NewPipeline(registry string) (*Pipeline, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, errors.New("AWS_REGION environment variable is required")
	}

	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	return &Pipeline{codebuild.New(sess), cloudwatchlogs.New(sess), registry}, nil
}

// Run starts an aws cloudbuild build operation
func (p *Pipeline) Run(ctx context.Context, logOut io.Writer, job *domain.BuildJob) error {
	projectName := aws.String(cbProject(job))
	projRes, err := p.cb.BatchGetProjects(&codebuild.BatchGetProjectsInput{
		Names: []*string{projectName},
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == codebuild.ErrCodeResourceNotFoundException {
			if err := p.createProject(job); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if len(projRes.Projects) == 0 {
		if err := p.createProject(job); err != nil {
			return err
		}
	}

	res, err := p.cb.StartBuild(&codebuild.StartBuildInput{
		ProjectName: projectName,
		EnvironmentVariablesOverride: []*codebuild.EnvironmentVariable{
			cbEnv("PULLR_REGISTRY", p.registry),
			cbEnv("PULLR_TAG", job.Tag),
			cbEnv("PULLR_OWNER", job.ImageOwner),
			cbEnv("PULLR_NAME", job.ImageName),
			cbEnv("PULLR_DOCKERFILE", job.Dockerfile),
			cbEnvSecret("PULLR_REGISTRY_USER"),
			cbEnvSecret("PULLR_REGISTRY_PASSWORD"),
		},
		SourceVersion: aws.String(job.CommitHash),
	})
	if err != nil {
		return err
	}

	numGetErr := 0
	var logs *codebuild.LogsLocation = nil
	for {
		res2, err := p.cb.BatchGetBuilds(&codebuild.BatchGetBuildsInput{Ids: []*string{res.Build.Id}})
		if err != nil {
			numGetErr++
			if numGetErr > 5 {
				return err
			}

			time.Sleep(time.Second * 10)
			continue
		} else if len(res2.Builds) != 1 {
			return errors.New("started build could not found in cloudbuild")
		}

		build := res2.Builds[0]
		if *build.BuildComplete {
			logs = build.Logs
			break
		}
	}

	// Allocate the buffer for at least one page of cloudwatch logs (1MB)
	logsInput := cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  logs.GroupName,
		LogStreamName: logs.StreamName,
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	p.logs.GetLogEventsPagesWithContext(timeoutCtx, &logsInput, func(output *cloudwatchlogs.GetLogEventsOutput, lastPage bool) bool {
		for _, ev := range output.Events {
			if _, err := io.WriteString(logOut, *ev.Message); err != nil {
				return false
			}
		}
		return true
	})
	cancel()

	return nil
}

func (p *Pipeline) createProject(job *domain.BuildJob) error {
	name := cbProject(job)
	sourceType := aws.String("")
	switch job.ImageRepo.Provider {
	case "github":
		sourceType = aws.String("GITHUB")
	default:
		return fmt.Errorf("unsupported source repository provider: %s", job.ImageRepo.Provider)
	}

	repoURL, err := job.ImageRepo.URL()
	if err != nil {
		return err
	}

	buildScriptOneLine := strings.Replace(buildScript, "\n", "", -1)
	buildSpec := fmt.Sprintf(buildSpecTemplate, buildScriptOneLine)

	_, err = p.cb.CreateProject(&codebuild.CreateProjectInput{
		Name: aws.String(name),
		Source: &codebuild.ProjectSource{
			Location:  aws.String(repoURL),
			Type:      sourceType,
			Buildspec: aws.String(buildSpec),
		},
		Environment: &codebuild.ProjectEnvironment{
			Type:           aws.String(codebuild.EnvironmentTypeLinuxContainer),
			Image:          aws.String("aws/codebuild/docker:17.09.0"),
			PrivilegedMode: aws.Bool(true),
			ComputeType:    aws.String(codebuild.ComputeTypeBuildGeneral1Small),
		},
		Artifacts: &codebuild.ProjectArtifacts{Type: aws.String(codebuild.ArtifactsTypeNoArtifacts)},
		Cache:     &codebuild.ProjectCache{Type: aws.String(codebuild.CacheTypeNoCache)},
	})
	return err
}

func cbEnv(key, value string) *codebuild.EnvironmentVariable {
	return &codebuild.EnvironmentVariable{
		Name:  aws.String(key),
		Value: aws.String(value),
		Type:  aws.String(codebuild.EnvironmentVariableTypePlaintext),
	}
}

func cbEnvSecret(key string) *codebuild.EnvironmentVariable {
	return &codebuild.EnvironmentVariable{
		Name:  aws.String(key),
		Value: aws.String(fmt.Sprintf("/CodeBuild/%s", key)),
		Type:  aws.String(codebuild.EnvironmentVariableTypeParameterStore),
	}
}

func cbProject(job *domain.BuildJob) string {
	// AWS CodeBuild only supports A-Za-z0-9\-\_ characters as project name
	normalizedKey := strings.Replace(job.ImageKey, ":", "_", -1)
	return fmt.Sprintf("pullr__%s__%s", job.ImageOwner, normalizedKey)
}
