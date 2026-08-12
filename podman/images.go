package podman

import (
	"fmt"

	"go.podman.io/podman/v6/pkg/bindings/images"
	"go.podman.io/podman/v6/pkg/domain/entities/types"
)

// PullImage downloads an image from a container registry.
func (connection *Connection) PullImage(image string) {
	_, err := images.Pull(connection.ctx, image, nil)
	if err != nil {
		fmt.Println(err)
	}
}

// GetImages retrieves a list of all images available in local storage.
func (connection *Connection) GetImages() []*types.ImageSummary {
	imagesList, err := images.List(connection.ctx, nil)
	if err != nil {
		fmt.Println(err)
	}
	return imagesList
}

// RemoveImages deletes one or more images from local storage.
func (connection *Connection) RemoveImages(imagesList []string) {
	_, err := images.Remove(connection.ctx, imagesList, nil)
	if err != nil {
		fmt.Println(err)
	}
}
