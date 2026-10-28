package podman

import (
	"fmt"

	"go.podman.io/podman/v6/pkg/bindings/images"
	"go.podman.io/podman/v6/pkg/domain/entities/types"
)

func (connection *Connection) PullImage(image string) {
	_, err := images.Pull(connection.ctx, image, nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Pulled Image", image)
}

func (connection *Connection) GetImages() []*types.ImageSummary {
	imagesList, err := images.List(connection.ctx, nil)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(imagesList)

	return imagesList
}

func (connection *Connection) RemoveImages(imagesList []string) {
	_, err := images.Remove(connection.ctx, imagesList, nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Removed Images", imagesList)
}

func (connection *Connection) CreateImage(dockerfile []string) *types.BuildReport {
	img, err := images.Build(connection.ctx, dockerfile, types.BuildOptions{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Created Image", img.ID)
	return img
}

func ImageSummary2StrArray(images []*types.ImageSummary) []string {
	strings := make([]string, len(images))
	for i, img := range images {
		strings[i] = img.ID
	}
	return strings
}
