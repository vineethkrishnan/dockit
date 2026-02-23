package dockerclient

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/vineethkrishnan/dockit/pkg/models"
)

// GetDiskSummary retrieves the global disk usage using Docker's /system/df endpoint.
// It maps the raw Docker SDK response into our domain model.
func (c *Client) GetDiskSummary(ctx context.Context) (*models.DiskSummary, error) {
	usage, err := c.api.DiskUsage(ctx, types.DiskUsageOptions{})
	if err != nil {
		return nil, err
	}

	summary := &models.DiskSummary{}

	// Calculate Images
	for _, img := range usage.Images {
		summary.Images.TotalCount++
		summary.Images.TotalSize += img.Size
		if img.Containers == 0 {
			summary.Images.Reclaimable += img.Size
		} else {
			summary.Images.ActiveCount++
		}
	}

	// Calculate Containers
	for _, container := range usage.Containers {
		summary.Containers.TotalCount++
		summary.Containers.TotalSize += container.SizeRw
		if container.State != "running" {
			summary.Containers.Reclaimable += container.SizeRw
		} else {
			summary.Containers.ActiveCount++
		}
	}

	// Calculate Volumes
	for _, vol := range usage.Volumes {
		summary.Volumes.TotalCount++
		// The Usage field is only populated if there is a daemon supporting it
		if vol.UsageData != nil {
			summary.Volumes.TotalSize += vol.UsageData.Size
			if vol.UsageData.RefCount == 0 {
				summary.Volumes.Reclaimable += vol.UsageData.Size
			} else {
				summary.Volumes.ActiveCount++
			}
		}
	}

	// Calculate Build Cache
	for _, cache := range usage.BuildCache {
		summary.BuildCache.TotalCount++
		summary.BuildCache.TotalSize += cache.Size
		if cache.InUse {
			summary.BuildCache.ActiveCount++
		} else {
			summary.BuildCache.Reclaimable += cache.Size
		}
	}

	// Calculate Totals
	summary.TotalSize = summary.Images.TotalSize + summary.Containers.TotalSize + summary.Volumes.TotalSize + summary.BuildCache.TotalSize
	summary.Reclaimable = summary.Images.Reclaimable + summary.Containers.Reclaimable + summary.Volumes.Reclaimable + summary.BuildCache.Reclaimable

	return summary, nil
}
