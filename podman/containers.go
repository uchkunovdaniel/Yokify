package podman

import (
	"context"
	"fmt"

	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/domain/entities/types"
	"go.podman.io/podman/v6/pkg/specgen"
)

// GetContainers retrieves a list of all containers from the Podman service.
// This includes both running and stopped containers.
func GetContainers(conn context.Context) []types.ListContainer {
	containersList, err := containers.List(conn, &containers.ListOptions{All: new(true)})
	if err != nil {
		fmt.Println(err)
	}

	return containersList
}

// RemoveContainer deletes a container identified by its name or ID.
// It is configured to ignore the error if the container does not exist.
func RemoveContainer(conn context.Context, container string) {
	_, err := containers.Remove(conn, container, &containers.RemoveOptions{
		Ignore: new(true),
	})
	if err != nil {
		fmt.Println(err)
	}
}

// StartContainer initiates the execution of a stopped container.
func StartContainer(conn context.Context, container string) {
	err := containers.Start(conn, container, nil)
	if err != nil {
		fmt.Println(err)
	}
}

// StopContainer halts a running container.
func StopContainer(conn context.Context, container string) {
	err := containers.Stop(conn, container, nil)
	if err != nil {
		fmt.Println(err)
	}
}

// RestartContainer stops and then starts a container.
func RestartContainer(conn context.Context, container string) {
	err := containers.Restart(conn, container, nil)
	if err != nil {
		fmt.Println(err)
	}
}

// CreateContainer creates a new container from the specified image.
// It uses default configuration settings for the new container.
func CreateContainer(conn context.Context, image string) {
	spec := specgen.NewSpecGenerator(image, false)
	_, err := containers.CreateWithSpec(conn, spec, nil)
	if err != nil {
		fmt.Println(err)
	}
}
