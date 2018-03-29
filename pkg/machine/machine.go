package machine

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/docker/machine/drivers/virtualbox"
	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/host"
	"github.com/mobingilabs/pullr/pkg/docker"
	"github.com/mobingilabs/pullr/pkg/domain"
)

const storePath = "/tmp/automatic"
const certsPath = "/tmp/automatic/certs"

// Machine builds docker images by provisioning docker machines
type Machine struct {
	client *libmachine.Client
	config *Config
}

// Config are the options for provisioning docker-machines
type Config struct {
	// Number of cpus
	CPU int
	// Amount of memory in megabytes
	Ram int
}

// ConfigFromMap transforms generic configuration into machine specific configuration
func ConfigFromMap(in map[string]string) (*Config, error) {
	cpuOpt, ok := in["cpu"]
	if !ok {
		cpuOpt = "1"
	}

	cpu, err := strconv.ParseInt(cpuOpt, 10, strconv.IntSize)
	if err != nil {
		return nil, err
	}

	ramOpt, ok := in["ram"]
	if !ok {
		ramOpt = "512"
	}

	ram, err := strconv.ParseInt(ramOpt, 10, strconv.IntSize)
	if err != nil {
		return nil, err
	}

	return &Config{int(cpu), int(ram)}, nil
}

// New creates a new machine instance. Machine provision docker-machines to
// build docker images
func New(config *Config) (*Machine, error) {
	client := libmachine.NewClient(storePath, certsPath)

	return &Machine{client, config}, nil
}

// Create, creates a new image builder by provisioning a new docker-machine
func (m *Machine) Create() (domain.ImageBuilder, error) {
	start := time.Now()
	hostname := fmt.Sprintf("builder%d", rand.NewSource(time.Now().UnixNano()).Int63()%10000)
	if !host.ValidateHostName(hostname) {
		return nil, fmt.Errorf("invalid hostname: %s", hostname)
	}

	driver := virtualbox.NewDriver(hostname, storePath)
	driver.CPU = m.config.CPU
	driver.Memory = m.config.Ram
	data, err := json.Marshal(driver)
	if err != nil {
		return nil, err
	}

	machineHost, err := m.client.NewHost("virtualbox", data)
	if err != nil {
		return nil, err
	}

	machineHost.HostOptions.EngineOptions.StorageDriver = "overlay"
	if err := m.client.Create(machineHost); err != nil {
		return nil, err
	}

	if err := m.client.Save(machineHost); err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stderr, "Provisioning took: %v", time.Since(start))

	dockerHost, err := machineHost.URL()
	if err != nil {
		return nil, err
	}

	d := docker.New(dockerHost, certsPath, true)
	return &Builder{m.client, machineHost, d}, nil
}

// Builder, builds and pushes docker images.
type Builder struct {
	client *libmachine.Client
	host   *host.Host
	d      *docker.Docker
}

// Close, destroys the docker-machine
func (b *Builder) Close() error {
	if err := b.host.Stop(); err != nil {
		return err
	}

	if err := b.host.Driver.Remove(); err != nil {
		return err
	}

	return b.client.Remove(b.host.Name)
}

// PushImage pushes a docker image to given docker registry
func (b *Builder) PushImage(ctx context.Context, out io.Writer, tag, registry, username, password string) error {
	return b.d.PushImage(ctx, out, tag, registry, username, password)
}

// BuildImage builds a docker image from given context and dockerfile
func (b *Builder) BuildImage(ctx context.Context, out io.Writer, ctxPath, dockerfile, tag string) error {
	return b.d.BuildImage(ctx, out, ctxPath, dockerfile, tag)
}
