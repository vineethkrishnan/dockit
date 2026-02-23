package dockerclient

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"

	"github.com/vineethkrishnan/dockit/pkg/models"
)

// GetContainers fetches all individual container metadata (both running and stopped)
func (c *Client) GetContainers(ctx context.Context) ([]*models.Container, error) {
	// Size: true is needed to ask Docker daemon to calculate SizeRw (which takes a bit longer).
	rawContainers, err := c.api.ContainerList(ctx, types.ContainerListOptions{All: true, Size: true})
	if err != nil {
		return nil, err
	}

	res := make([]*models.Container, 0, len(rawContainers))
	for _, cnt := range rawContainers {
		m := &models.Container{
			ID:       cnt.ID,
			Names:    cnt.Names,
			State:    cnt.State,
			Status:   cnt.Status,
			Image:    cnt.Image,
			ImageID:  cnt.ImageID,
			SizeRw:   cnt.SizeRw,
			SizeRoot: cnt.SizeRootFs,
			Created:  time.Unix(cnt.Created, 0),
		}
		res = append(res, m)
	}
	return res, nil
}

// GetImages fetches all individual image metadata
func (c *Client) GetImages(ctx context.Context) ([]*models.Image, error) {
	rawImages, err := c.api.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		return nil, err
	}

	res := make([]*models.Image, 0, len(rawImages))
	for _, img := range rawImages {
		dangling := false
		// Docker represents dangling images by setting the RepoTags to <none>:<none> usually, 
		// but the SDK uses the labels or simply having empty tags.
		if len(img.RepoTags) == 0 || (len(img.RepoTags) == 1 && img.RepoTags[0] == "<none>:<none>") {
			dangling = true
		}

		m := &models.Image{
			ID:         img.ID,
			RepoTags:   img.RepoTags,
			Size:       img.Size,
			Created:    time.Unix(img.Created, 0),
			Containers: img.Containers, // -1 means it wasn't populated in ImageList, requires /system/df or detailed inspect. We'll correlate manually.
			Dangling:   dangling,
		}
		res = append(res, m)
	}
	return res, nil
}

// GetVolumes fetches all individual volume metadata
func (c *Client) GetVolumes(ctx context.Context) ([]*models.Volume, error) {
	rawVolumes, err := c.api.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, err
	}

	res := make([]*models.Volume, 0, len(rawVolumes.Volumes))
	for _, vol := range rawVolumes.Volumes {
		vt := time.Time{}
		if vol.CreatedAt != "" {
			parsed, err := time.Parse(time.RFC3339, vol.CreatedAt)
			if err == nil {
				vt = parsed
			}
		}

		var size, usage int64
		if vol.UsageData != nil {
			size = vol.UsageData.Size
			usage = vol.UsageData.RefCount
		}

		m := &models.Volume{
			Name:       vol.Name,
			Driver:     vol.Driver,
			Mountpoint: vol.Mountpoint,
			CreatedAt:  vt,
			Size:       size,
			UsageCount: usage, // Depending on Docker version, this might be -1 if /system/df wasn't called. We'll correlate manually in Analyzer.
		}
		res = append(res, m)
	}
	return res, nil
}
