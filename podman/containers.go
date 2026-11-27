package podman

import (
	"context"
	"errors"
	"fmt"

	nettypes "go.podman.io/common/libnetwork/types"
	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/domain/entities/types"
	"go.podman.io/podman/v6/pkg/specgen"
)

// PortMapping publishes a container port on the host.
type PortMapping struct {
	HostIP        string // "" = all interfaces; use "127.0.0.1" for this machine only
	HostPort      uint16 // 0 = Podman picks a free port
	ContainerPort uint16 // required
	Protocol      string // "tcp" (default), "udp" or "sctp"
}

// CreateContainerOptions configures a new container. Every field is optional.
type CreateContainerOptions struct {
	Name  string            // "" = Podman generates one
	Ports []PortMapping     // published ports
	Env   map[string]string // environment variables
	Start bool              // start the container after creating it
}

// GetContainers includes both running and stopped containers.
func (connection *Connection) GetContainers(ctx context.Context) ([]types.ListContainer, error) {
	ctx, done := connection.call(ctx)
	defer done()
	return containers.List(ctx, &containers.ListOptions{All: new(true)})
}

// RemoveContainer deletes a container identified by its name or ID.
func (connection *Connection) RemoveContainer(ctx context.Context, container string, force bool) error {
	ctx, done := connection.call(ctx)
	defer done()

	rmReports, err := containers.Remove(ctx, container, &containers.RemoveOptions{
		Force: &force,
	})
	if err != nil {
		return fmt.Errorf("remove container %s: %w", container, err)
	}

	var errs []error
	for _, r := range rmReports {
		if r.Err != nil {
			errs = append(errs, fmt.Errorf("container %s: %w", r.Id, r.Err))
		}
	}
	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("remove container %s: %w", container, err)
	}
	return nil
}

// StartContainer starts a non-running container.
// The argument can be a container name or a partial/full ID.
func (connection *Connection) StartContainer(ctx context.Context, container string) error {
	ctx, done := connection.call(ctx)
	defer done()
	return containers.Start(ctx, container, nil)
}

// StopContainer stops a running container.
// The argument can be a container name or a partial/full ID.
func (connection *Connection) StopContainer(ctx context.Context, container string) error {
	ctx, done := connection.call(ctx)
	defer done()

	return containers.Stop(ctx, container, nil)
}

// RestartContainer restarts a running container.
// The argument can be a container name or a partial/full ID.
func (connection *Connection) RestartContainer(ctx context.Context, container string) error {
	ctx, done := connection.call(ctx)
	defer done()
	return containers.Restart(ctx, container, nil)
}

// CreateContainerFromImage creates a new container from an existing image.
// It takes the image name or ID as an argument and returns a ContainerCreateResponse containing details about the newly created container.
func (connection *Connection) CreateContainerFromImage(ctx context.Context, image string,
	options CreateContainerOptions,
) (types.ContainerCreateResponse, error) {
	ctx, done := connection.call(ctx)
	defer done()

	spec := specgen.NewSpecGenerator(image, false)
	spec.Name = options.Name
	spec.Env = options.Env
	for _, p := range options.Ports {
		if p.ContainerPort == 0 {
			return types.ContainerCreateResponse{}, fmt.Errorf("port mapping %+v: container port is required", p)
		}
		spec.PortMappings = append(spec.PortMappings, nettypes.PortMapping{
			HostIP:        p.HostIP,
			HostPort:      p.HostPort,
			ContainerPort: p.ContainerPort,
			Protocol:      p.Protocol,
		})
	}

	resp, err := containers.CreateWithSpec(ctx, spec, nil)
	if err != nil {
		return resp, fmt.Errorf("create container from %s: %w", image, err)
	}
	if options.Start {
		if errStart := containers.Start(ctx, resp.ID, nil); errStart != nil {
			return resp, fmt.Errorf("start container %s: %w", resp.ID, errStart)
		}
	}
	return resp, nil
}

// CreateContainerFromFile creates a new container from a Dockerfile/Containerfile.
// It takes the project root directory as an argument and returns a
// ContainerCreateResponse containing details about the newly created container.
func (connection *Connection) CreateContainerFromFile(ctx context.Context, projectRoot string, name string,
	containerOptions CreateContainerOptions, imageOptions CreateImageOptions,
) (types.ContainerCreateResponse, error) {
	img, err := connection.CreateImage(ctx, projectRoot, name, imageOptions)
	if err != nil {
		return types.ContainerCreateResponse{}, err
	}
	return connection.CreateContainerFromImage(ctx, img.ID, containerOptions)
}
