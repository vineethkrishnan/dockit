package analyzer

import (
	"github.com/vineethkrishnan/dockit/pkg/models"
)

// CorrelatedData holds the cross-referenced resources.
type CorrelatedData struct {
	Containers []*models.Container
	Images     []*models.Image
	Volumes    []*models.Volume
}

// Analyze processes the raw lists of Docker resources and maps their relationships.
func Analyze(containers []*models.Container, images []*models.Image, volumes []*models.Volume) *CorrelatedData {
	// 1. Manually calculate how many containers are using each image.
	// This ensures that even if the Docker daemon didn't populate Image.Containers in ImageList,
	// we have a localized true count based on actual local containers.
	imgUsage := make(map[string]int64)
	for _, c := range containers {
		imgUsage[c.ImageID]++
	}

	for _, img := range images {
		// If Docker natively reported a value (and it's not -1), keep the max of natively reported or locally observed
		localCount := imgUsage[img.ID]
		if localCount > img.Containers {
			img.Containers = localCount
		}
	}

	// Future: Manually correlate container volume mounts to volumes if usage_count is missing.

	return &CorrelatedData{
		Containers: containers,
		Images:     images,
		Volumes:    volumes,
	}
}
