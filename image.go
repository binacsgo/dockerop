package dockerop

import (
	"bytes"
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// ImageList list all the images
func (c *OpClient) ImageList(ctx context.Context) ([]types.ImageSummary, error) {
	return c.client.ImageList(context.Background(), types.ImageListOptions{})
}

// ImagePull pull the image, options has the cred of registry.
// ref such as "docker.io/library/nginx:latest"
// return `docker pull stdout information` and `error`
func (c *OpClient) ImagePull(ctx context.Context, ref string, options types.ImagePullOptions) (string, error) {
	reader, err := c.client.ImagePull(ctx, ref, options)
	if err != nil {
		return "", err
	}
	info := new(bytes.Buffer)
	info.ReadFrom(reader)
	return info.String(), nil
}

// ImageRemove remove the image
func (c *OpClient) ImageRemove(ctx context.Context, imageID string) error {
	// items
	_, err := c.client.ImageRemove(ctx, imageID, types.ImageRemoveOptions{Force: true, PruneChildren: true})
	return err
}

// ImagesPrune prune the image
func (c *OpClient) ImagesPrune(ctx context.Context) error {
	_, err := c.client.ImagesPrune(ctx, filters.Args{})
	return err
}
