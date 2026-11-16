package podman

import (
	"context"
	"errors"
	"fmt"

	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/domain/entities/types"
	"go.podman.io/podman/v6/pkg/specgen"
)

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
func (connection *Connection) CreateContainerFromImage(ctx context.Context, image string) (types.ContainerCreateResponse, error) {
	ctx, done := connection.call(ctx)
	defer done()

	spec := specgen.NewSpecGenerator(image, false)
	return containers.CreateWithSpec(ctx, spec, nil)
}

// CreateContainerFromFile creates a new container from a Dockerfile/Containerfile.
// It takes the project root directory as an argument and returns a
// ContainerCreateResponse containing details about the newly created container.
func (connection *Connection) CreateContainerFromFile(
	ctx context.Context, projectRoot string, name string, options CreateImageOptions,
) (types.ContainerCreateResponse, error) {
	img, err := connection.CreateImage(ctx, projectRoot, name, options)
	if err != nil {
		return types.ContainerCreateResponse{}, err
	}
	return connection.CreateContainerFromImage(ctx, img.ID)
}
