package dockerclient

import (
	"context"

	"github.com/docker/docker/api/types"
)

// RemoveContainer permanently deletes a container by ID.
func (c *Client) RemoveContainer(ctx context.Context, id string) error {
	return c.api.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})
}

// RemoveImage permanently deletes an image by ID.
func (c *Client) RemoveImage(ctx context.Context, id string) error {
	_, err := c.api.ImageRemove(ctx, id, types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
	return err
}

// RemoveVolume permanently deletes a volume by name.
func (c *Client) RemoveVolume(ctx context.Context, name string) error {
	return c.api.VolumeRemove(ctx, name, true)
}
