package dockerclient

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/vineethkrishnan/dockit/pkg/models"
)

// GetLogMetrics retrieves the log file paths for all containers and stats them on disk.
func (c *Client) GetLogMetrics(ctx context.Context) ([]*models.LogMetrics, error) {
	// 1. Get all containers
	containers, err := c.api.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	metrics := make([]*models.LogMetrics, 0, len(containers))

	for _, cnt := range containers {
		// 2. Inspect each container to find exactly where its log is stored on the host
		info, err := c.api.ContainerInspect(ctx, cnt.ID)
		if err != nil {
			continue // Skip containers that error out during inspection
		}

		name := cnt.ID[:10]
		if len(cnt.Names) > 0 {
			name = cnt.Names[0]
		}
		
		m := &models.LogMetrics{
			ContainerID:   cnt.ID,
			ContainerName: name,
			LogPath:       info.LogPath,
			HasLogDriver:  info.HostConfig != nil && info.HostConfig.LogConfig.Type != "",
		}

		// 3. Stat the actual file on disk if it exists
		if info.LogPath != "" {
			stat, err := os.Stat(info.LogPath)
			if err == nil {
				m.LogSize = stat.Size()
			}
		}

		metrics = append(metrics, m)
	}

	return metrics, nil
}
