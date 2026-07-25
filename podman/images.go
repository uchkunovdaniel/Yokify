package podman

import (
	"context"
	"fmt"

	"go.podman.io/podman/v6/pkg/bindings/images"
	"go.podman.io/podman/v6/pkg/domain/entities/types"
)

// PullImage downloads an image from a container registry.
func PullImage(conn context.Context, image string) {
	_, err := images.Pull(conn, image, nil)
	if err != nil {
		fmt.Println(err)
	}
}

// GetImages retrieves a list of all images available in local storage.
func GetImages(conn context.Context) []*types.ImageSummary {
	imagesList, err := images.List(conn, nil)
	if err != nil {
		fmt.Println(err)
	}
	return imagesList
}

// RemoveImages deletes one or more images from local storage.
func RemoveImages(conn context.Context, imagesList []string) {
	_, err := images.Remove(conn, imagesList, nil)
	if err != nil {
		fmt.Println(err)
	}
}
