package podman

import (
	"fmt"

	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/domain/entities/types"
	"go.podman.io/podman/v6/pkg/specgen"
)

// GetContainers retrieves a list of all containers from the Podman service.
// This includes both running and stopped containers.
func (connection *Connection) GetContainers() []types.ListContainer {
	containersList, err := containers.List(connection.ctx, &containers.ListOptions{All: new(true)})
	if err != nil {
		fmt.Println(err)
	}

	return containersList
}

// RemoveContainer deletes a container identified by its name or ID.
// It is configured to ignore the error if the container does not exist.
func (connection *Connection) RemoveContainer(container string) {
	_, err := containers.Remove(connection.ctx, container, &containers.RemoveOptions{
		Ignore: new(true),
	})
	if err != nil {
		fmt.Println(err)
	}
}

// StartContainer initiates the execution of a stopped container.
func (connection *Connection) StartContainer(container string) {
	err := containers.Start(connection.ctx, container, nil)
	if err != nil {
		fmt.Println(err)
	}
}

// StopContainer halts a running container.
func (connection *Connection) StopContainer(container string) {
	err := containers.Stop(connection.ctx, container, nil)
	if err != nil {
		fmt.Println(err)
	}
}

// RestartContainer stops and then starts a container.
func (connection *Connection) RestartContainer(container string) {
	err := containers.Restart(connection.ctx, container, nil)
	if err != nil {
		fmt.Println(err)
	}
}

// CreateContainer creates a new container from the specified image.
// It uses default configuration settings for the new container.
func (connection *Connection) CreateContainer(image string) {
	spec := specgen.NewSpecGenerator(image, false)
	_, err := containers.CreateWithSpec(connection.ctx, spec, nil)
	if err != nil {
		fmt.Println(err)
	}
}
