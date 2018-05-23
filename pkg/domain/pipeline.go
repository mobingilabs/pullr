package domain

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// PipelineConfig is the configuration data for build pipeline
type PipelineConfig struct {
	CloneDir         string
	RegistryURL      string
	RegistryUser     string
	RegistryPassword string
}

// RepositoryCloner clones source code
type RepositoryCloner interface {
	// CloneRepository clones the given source repository into target directory and reports
	// back the path for cloned repository.
	CloneRepository(ctx context.Context, out io.Writer, target string, repo SourceRepository, username, token string) error
}

// ImageBuilderFactory creates ImageBuilders. Each running pipeline gets its own image
// builder. Image builders provides an opportunity cleanup the resources created during
// the build.
type ImageBuilderFactory interface {
	Create() (ImageBuilder, error)
}

// ImageBuilder builds container images. When it is closed, it cleans up the
// resources created during building an image.
type ImageBuilder interface {
	io.Closer
	// BuildImage builds a container image from a Dockerfile located at ctxPath.
	// ctxPath is also used as build context
	BuildImage(ctx context.Context, out io.Writer, ctxPath, dockerfile, tag string) error

	// PushImage pushes a container image to the given registry
	PushImage(ctx context.Context, out io.Writer, tag, registry, username, password string) error
}

type Pipeline interface {
	Run(ctx context.Context, logOut io.Writer, job *BuildJob) error
}

// HostedPipeline is the image build pipeline
type HostedPipeline struct {
	config         PipelineConfig
	logger         Logger
	cloners        map[string]RepositoryCloner
	builderFactory ImageBuilderFactory
	randSource     rand.Source
}

// NewPipeline creates a build pipeline for given job
func NewPipeline(config PipelineConfig, logger Logger, cloners map[string]RepositoryCloner, builderFactory ImageBuilderFactory) *HostedPipeline {
	randSource := rand.NewSource(time.Now().UnixNano())
	return &HostedPipeline{config, logger, cloners, builderFactory, randSource}
}

// Run, runs the build pipeline against the given job
func (p *HostedPipeline) Run(ctx context.Context, out io.Writer, job *BuildJob) (err error) {
	cloner, ok := p.cloners[job.ImageRepo.Provider]
	if !ok {
		return ErrSourceUnsupportedProvider
	}

	dirname := fmt.Sprintf("%s_%d", job.ImageRepo.Name, p.randSource.Int63())
	dir := filepath.Join(p.config.CloneDir, dirname)
	defer os.RemoveAll(dir)

	err = cloner.CloneRepository(ctx, out, dir, job.ImageRepo, job.VcsUsername, job.VcsToken)
	if err != nil {
		return fmt.Errorf("pipeline: clone: %v", err)
	}

	builder, err := p.builderFactory.Create()
	if err != nil {
		return err
	}
	defer builder.Close()
	tag := fmt.Sprintf("%s/%s:%s", job.ImageOwner, job.ImageName, job.Tag)
	err = builder.BuildImage(ctx, out, dir, job.Dockerfile, tag)
	if err != nil {
		return fmt.Errorf("pipeline: build: %v", err)
	}

	err = builder.PushImage(ctx, out, tag, p.config.RegistryURL, p.config.RegistryUser, p.config.RegistryPassword)
	if err != nil {
		return fmt.Errorf("pipeline: push: %v", err)
	}
	if err := os.RemoveAll(dir); err != nil {
		p.logger.Errorf("pipeline: remove repo dir: %v", err)
	}

	return nil
}
