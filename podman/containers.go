package podman

import (
	"errors"
	"fmt"

	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/domain/entities/types"
	"go.podman.io/podman/v6/pkg/specgen"
)

// GetContainers includes both running and stopped containers.
func (connection *Connection) GetContainers() ([]types.ListContainer, error) {
	return containers.List(connection.ctx, &containers.ListOptions{All: new(true)})
}

// RemoveContainer deletes a container identified by its name or ID.
func (connection *Connection) RemoveContainer(container string, force bool) error {
	rmReports, err := containers.Remove(connection.ctx, container, &containers.RemoveOptions{
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
func (connection *Connection) StartContainer(container string) error {
	return containers.Start(connection.ctx, container, nil)
}

// StopContainer stops a running container.
// The argument can be a container name or a partial/full ID.
func (connection *Connection) StopContainer(container string) error {
	return containers.Stop(connection.ctx, container, nil)
}

// RestartContainer restarts a running container.
// The argument can be a container name or a partial/full ID.
func (connection *Connection) RestartContainer(container string) error {
	return containers.Restart(connection.ctx, container, nil)
}

// CreateContainerFromImage creates a new container from an existing image.
// It takes the image name or ID as an argument and returns a ContainerCreateResponse containing details about the newly created container.
func (connection *Connection) CreateContainerFromImage(image string) (types.ContainerCreateResponse, error) {
	spec := specgen.NewSpecGenerator(image, false)
	return containers.CreateWithSpec(connection.ctx, spec, nil)
}

// CreateContainerFromFile creates a new container from a Dockerfile.
// It takes the project root directory as an argument and returns a ContainerCreateResponse containing details about the newly created container.
func (connection *Connection) CreateContainerFromFile(projectRoot string) (types.ContainerCreateResponse, error) {
	img, err := connection.CreateImage(projectRoot)
	if err != nil {
		return types.ContainerCreateResponse{}, err
	}
	return connection.CreateContainerFromImage(img.ID)
}
