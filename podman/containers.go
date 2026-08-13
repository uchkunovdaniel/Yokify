package podman

import (
	"fmt"

	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/domain/entities/types"
	"go.podman.io/podman/v6/pkg/specgen"
)

// GetContainers includes both running and stopped containers.
func (connection *Connection) GetContainers() []types.ListContainer {
	containersList, err := containers.List(connection.ctx, &containers.ListOptions{All: new(true)})
	if err != nil {
		fmt.Println(err)
	}

	return containersList
}

// RemoveContainer deletes a container identified by its name or ID.
func (connection *Connection) RemoveContainer(container string) {
	_, err := containers.Remove(connection.ctx, container, &containers.RemoveOptions{
		Ignore: new(true),
	})
	if err != nil {
		fmt.Println(err)
	}
}

func (connection *Connection) StartContainer(container string) {
	err := containers.Start(connection.ctx, container, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (connection *Connection) StopContainer(container string) {
	err := containers.Stop(connection.ctx, container, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (connection *Connection) RestartContainer(container string) {
	err := containers.Restart(connection.ctx, container, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (connection *Connection) CreateContainerFromImage(image string) types.ContainerCreateResponse {
	spec := specgen.NewSpecGenerator(image, false)
	container, err := containers.CreateWithSpec(connection.ctx, spec, nil)
	if err != nil {
		fmt.Println(err)
	}
	return container
}

func (connection *Connection) CreateContainerFromFile(dockerfile []string) types.ContainerCreateResponse {
	img := connection.CreateImage(dockerfile)
	return connection.CreateContainerFromImage(img.ID)
}
