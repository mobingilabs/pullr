package docker

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/mobingilabs/pullr/pkg/domain"
)

type Factory struct {
	host     string
	certPath string
	tls      bool
}

func NewFactory(host, certPath string) *Factory {
	tls := certPath != ""
	return &Factory{host, certPath, tls}
}

func (f *Factory) Create() (domain.ImageBuilder, error) {
	return New(f.host, f.certPath, f.tls), nil
}

// Docker wraps docker executable to build and push images
type Docker struct {
	env []string
}

// New creates a new docker image builder with given docker host and certs path. Use empty
// strings for both host and certPath for using the default docker configuration
func New(host string, certPath string, tls bool) *Docker {
	var env []string
	if host != "" {
		env = append(env, fmt.Sprintf("DOCKER_HOST=%s", host))
	}
	if certPath != "" {
		env = append(env, fmt.Sprintf("DOCKER_CERT_PATH=%s", certPath))
	}
	if tls {
		env = append(env, "DOCKER_TLS_VERIFY=1")
	}

	return &Docker{env}
}

func (d *Docker) Close() error {
	return nil
}

// TagImage tags and image
func (d *Docker) TagImage(ctx context.Context, out io.Writer, imageTag, newTag string) error {
	cmd := exec.CommandContext(ctx, "docker", "tag", imageTag, newTag)
	cmd.Env = d.env
	cmd.Stdout = out
	cmd.Stderr = out
	return cmd.Run()
}

// Login, logs in to a docker registry
func (d *Docker) Login(ctx context.Context, out io.Writer, registry, username, password string) error {
	cmd := exec.CommandContext(ctx, "docker", "login", "-u", username, "-p", password, registry)
	cmd.Env = d.env
	cmd.Stdout = out
	cmd.Stderr = out
	return cmd.Run()
}

// PushImage pushes a container image to the given registry
func (d *Docker) PushImage(ctx context.Context, out io.Writer, tag, registry, username, password string) error {
	if err := d.Login(ctx, out, registry, username, password); err != nil {
		return err
	}

	remoteTag := fmt.Sprintf("%s/%s", registry, tag)
	if err := d.TagImage(ctx, out, tag, remoteTag); err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "docker", "push", remoteTag)
	cmd.Env = d.env
	cmd.Stderr = out
	cmd.Stdout = out
	return cmd.Run()
}

// BuildImage builds a container image from a Dockerfile located at ctxPath.
// ctxPath is also used as build context
func (d *Docker) BuildImage(ctx context.Context, out io.Writer, ctxPath, dockerfile, tag string) error {
	cmd := exec.Command("docker", "build", "-t", tag, "-f", dockerfile, ".")
	cmd.Dir = ctxPath
	cmd.Env = d.env
	cmd.Stderr = out
	cmd.Stdout = out

	return cmd.Run()
}
