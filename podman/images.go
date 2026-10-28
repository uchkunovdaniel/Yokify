package podman

import (
	"fmt"
	"path"

	"go.podman.io/podman/v6/pkg/bindings/images"
	"go.podman.io/podman/v6/pkg/domain/entities/types"
)

// PullImage pulls an image from a registry
func (connection *Connection) PullImage(image string) {
	_, err := images.Pull(connection.ctx, image, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Pulled Image", image)
}

// GetImages returns a list of images in local storage.
func (connection *Connection) GetImages() []*types.ImageSummary {
	imagesList, err := images.List(connection.ctx, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println(imagesList)

	return imagesList
}

// RemoveImages removes one or more images from the local storage.
func (connection *Connection) RemoveImages(imagesList []string) {
	_, err := images.Remove(connection.ctx, imagesList, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Removed Images", imagesList)
}

// CreateImage performs a container image build on the remote API using the standard build API.
// Prepares build contexts and container files by creating tar archives from local directories, processes build secrets and authentication, and streams the build to the remote server.
// Returns a BuildReport containing the final image ID and save format.
func (connection *Connection) CreateImage(projectRoot string) *types.BuildReport {
	containerfile := path.Join(projectRoot, "Containerfile")
	img, err := images.Build(connection.ctx, []string{containerfile}, types.BuildOptions{
		ContextDirectory: projectRoot,
	})
	if err != nil {
		fmt.Println(err)
		return &types.BuildReport{}
	}
	fmt.Println("Created Image", img.ID)
	return img
}

// ImageSummary2StrSlice converts a slice of ImageSummary objects to a slice of strings containing the image IDs.
func ImageSummary2StrSlice(images []*types.ImageSummary) []string {
	strings := make([]string, len(images))
	for i, img := range images {
		strings[i] = img.ID
	}
	return strings
}
