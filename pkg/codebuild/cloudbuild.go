package codebuild

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
	awscb "github.com/aws/aws-sdk-go/service/codebuild"
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

// Pipeline builds and pushes docker images to registries by using aws codebuild service
type Pipeline struct {
	cb       *awscb.CodeBuild
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

	return &Pipeline{awscb.New(sess), cloudwatchlogs.New(sess), registry}, nil
}

// Run starts an aws codebuild build operation
func (p *Pipeline) Run(ctx context.Context, logOut io.Writer, job *domain.BuildJob) (domain.BuildStatus, error) {
	projectName := aws.String(cbProject(job))
	projRes, err := p.cb.BatchGetProjects(&awscb.BatchGetProjectsInput{
		Names: []*string{projectName},
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == awscb.ErrCodeResourceNotFoundException {
			if err := p.createProject(job); err != nil {
				return domain.BuildFailed, err
			}
		} else {
			return domain.BuildFailed, err
		}
	}
	if len(projRes.Projects) == 0 {
		if err := p.createProject(job); err != nil {
			return domain.BuildFailed, err
		}
	}

	res, err := p.cb.StartBuild(&awscb.StartBuildInput{
		ProjectName: projectName,
		EnvironmentVariablesOverride: []*awscb.EnvironmentVariable{
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
		return domain.BuildFailed, err
	}

	numGetErr := 0
	status := domain.BuildInProgress
	var logs *awscb.LogsLocation = nil
	for {
		res2, err := p.cb.BatchGetBuilds(&awscb.BatchGetBuildsInput{Ids: []*string{res.Build.Id}})
		if err != nil {
			numGetErr++
			if numGetErr > 5 {
				return domain.BuildFailed, err
			}

			time.Sleep(time.Second * 10)
			continue
		} else if len(res2.Builds) != 1 {
			return domain.BuildFailed, errors.New("started build could not found in codebuild")
		}

		build := res2.Builds[0]
		if *build.BuildComplete {
			logs = build.Logs
			switch *build.BuildStatus {
			case awscb.StatusTypeSucceeded:
				status = domain.BuildSucceed
			case awscb.StatusTypeFailed:
				status = domain.BuildFailed
			case awscb.StatusTypeTimedOut:
				status = domain.BuildTimeout
			default:
				// TODO: should we mark it as in progress if we don't know the status
				status = domain.BuildInProgress
			}
			break
		}
	}

	if logs != nil {
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
	}

	return status, nil
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

	_, err = p.cb.CreateProject(&awscb.CreateProjectInput{
		Name: aws.String(name),
		Source: &awscb.ProjectSource{
			Location:  aws.String(repoURL),
			Type:      sourceType,
			Buildspec: aws.String(buildSpec),
		},
		Environment: &awscb.ProjectEnvironment{
			Type:           aws.String(awscb.EnvironmentTypeLinuxContainer),
			Image:          aws.String("aws/codebuild/docker:17.09.0"),
			PrivilegedMode: aws.Bool(true),
			ComputeType:    aws.String(awscb.ComputeTypeBuildGeneral1Small),
		},
		Artifacts:   &awscb.ProjectArtifacts{Type: aws.String(awscb.ArtifactsTypeNoArtifacts)},
		Cache:       &awscb.ProjectCache{Type: aws.String(awscb.CacheTypeNoCache)},
		ServiceRole: aws.String(os.Getenv("PULLR_CODEBUILD_SERVICE_ROLE")),
	})
	return err
}

func cbEnv(key, value string) *awscb.EnvironmentVariable {
	return &awscb.EnvironmentVariable{
		Name:  aws.String(key),
		Value: aws.String(value),
		Type:  aws.String(awscb.EnvironmentVariableTypePlaintext),
	}
}

func cbEnvSecret(key string) *awscb.EnvironmentVariable {
	return &awscb.EnvironmentVariable{
		Name:  aws.String(key),
		Value: aws.String(fmt.Sprintf("/CodeBuild/%s", key)),
		Type:  aws.String(awscb.EnvironmentVariableTypeParameterStore),
	}
}

func cbProject(job *domain.BuildJob) string {
	// AWS CodeBuild only supports A-Za-z0-9\-\_ characters as project name
	normalizedKey := strings.Replace(job.ImageKey, ":", "_", -1)
	return fmt.Sprintf("pullr__%s__%s", job.ImageOwner, normalizedKey)
}
