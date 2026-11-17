// Package podman wraps the Podman API for Yōkify.
package podman

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"go.podman.io/podman/v6/pkg/bindings/images"
	"go.podman.io/podman/v6/pkg/domain/entities/types"
)

// CreateImageOptions defines options for creating an image from a Dockerfile/Containerfile.
type CreateImageOptions struct {
	Containerfile string
	Out           io.Writer
}

// PullImage pulls an image from a registry
func (connection *Connection) PullImage(ctx context.Context, image string) ([]string, error) {
	ctx, done := connection.call(ctx)
	defer done()

	return images.Pull(ctx, image, nil)
}

// GetImages returns a list of images in local storage.
func (connection *Connection) GetImages(ctx context.Context) ([]*types.ImageSummary, error) {
	ctx, done := connection.call(ctx)
	defer done()

	return images.List(ctx, nil)
}

// RemoveImages removes one or more images from the local storage.
func (connection *Connection) RemoveImages(ctx context.Context, imagesList []string) error {
	ctx, done := connection.call(ctx)
	defer done()

	_, errs := images.Remove(ctx, imagesList, nil)
	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("remove images %v: %w", imagesList, err)
	}
	return nil
}

// CreateImage performs a container image build on the remote API using the standard build API.
// Prepares build contexts and container files by creating tar archives from local directories,
// processes build secrets and authentication, and streams the build to the remote server.
// Returns a BuildReport containing the final image ID and save format.
func (connection *Connection) CreateImage(
	ctx context.Context, projectRoot string, name string, options CreateImageOptions,
) (*types.BuildReport, error) {
	ctx, done := connection.call(ctx)
	defer done()

	file := "Containerfile"
	if options.Containerfile != "" {
		file = options.Containerfile
	}
	if name == "" {
		abs, err := filepath.Abs(projectRoot)
		if err != nil {
			return nil, fmt.Errorf("resolve project path %q: %w", projectRoot, err)
		}
		name = filepath.Base(abs)
	}
	containerfile := filepath.Join(projectRoot, file)
	return images.Build(ctx, []string{containerfile}, types.BuildOptions{
		ContextDirectory: projectRoot,
		Output:           strings.ToLower(name),
		Out:              options.Out,
	})
}
