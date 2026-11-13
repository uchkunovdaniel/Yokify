package podman

import (
	"errors"
	"fmt"
	"path/filepath"

	"go.podman.io/buildah/define"
	"go.podman.io/podman/v6/pkg/bindings/images"
	"go.podman.io/podman/v6/pkg/domain/entities/types"
)

type CreateImageOptions struct {
	Containerfile string `default:"Containerfile"`
}

// PullImage pulls an image from a registry
func (connection *Connection) PullImage(image string) ([]string, error) {
	return images.Pull(connection.ctx, image, nil)
}

// GetImages returns a list of images in local storage.
func (connection *Connection) GetImages() ([]*types.ImageSummary, error) {
	return images.List(connection.ctx, nil)
}

// RemoveImages removes one or more images from the local storage.
func (connection *Connection) RemoveImages(imagesList []string) error {
	_, errs := images.Remove(connection.ctx, imagesList, nil)
	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("remove images %v: %w", imagesList, err)
	}
	return nil
}

// CreateImage performs a container image build on the remote API using the standard build API.
// Prepares build contexts and container files by creating tar archives from local directories,
// processes build secrets and authentication, and streams the build to the remote server.
// Returns a BuildReport containing the final image ID and save format.
func (connection *Connection) CreateImage(projectRoot string, name string, options *CreateImageOptions) (*types.BuildReport, error) {
	if options == nil {
		options = &CreateImageOptions{}
		options.Containerfile = "Containerfile"
	}

	containerfile := filepath.Join(projectRoot, options.Containerfile)
	return images.Build(connection.ctx, []string{containerfile}, types.BuildOptions{
		BuildOptions: define.BuildOptions{
			ContextDirectory: projectRoot,
			Output:           name,
		},
	})
}
