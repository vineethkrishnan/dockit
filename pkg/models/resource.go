package models

import "time"

// DiskSummary represents the global state of Docker disk usage.
type DiskSummary struct {
	Images       ResourceSummary `json:"images"`
	Containers   ResourceSummary `json:"containers"`
	Volumes      ResourceSummary `json:"volumes"`
	BuildCache   ResourceSummary `json:"build_cache"`
	TotalSize    int64           `json:"total_size"`
	Reclaimable  int64           `json:"reclaimable"`
}

// ResourceSummary represents the aggregated summary for a specific resource type.
type ResourceSummary struct {
	TotalCount  int   `json:"total_count"`
	ActiveCount int   `json:"active_count"`
	TotalSize   int64 `json:"total_size"`
	Reclaimable int64 `json:"reclaimable"`
}

// Image represents a parsed Docker image.
type Image struct {
	ID          string    `json:"id"`
	RepoTags    []string  `json:"repo_tags"`
	Size        int64     `json:"size"`
	Created     time.Time `json:"created"`
	Containers  int64     `json:"containers"` // Number of containers using this image
	Dangling    bool      `json:"dangling"`
}

// Container represents a parsed Docker container.
type Container struct {
	ID        string    `json:"id"`
	Names     []string  `json:"names"`
	State     string    `json:"state"` // running, exited, etc.
	Status    string    `json:"status"`
	Image     string    `json:"image"`
	ImageID   string    `json:"image_id"`
	SizeRw    int64     `json:"size_rw"` // Size of writable layer
	SizeRoot  int64     `json:"size_rootfs"`
	Created   time.Time `json:"created"`
}

// Volume represents a parsed Docker volume.
type Volume struct {
	Name       string    `json:"name"`
	Driver     string    `json:"driver"`
	Mountpoint string    `json:"mountpoint"`
	CreatedAt  time.Time `json:"created_at"`
	Size       int64     `json:"size"`
	UsageCount int64     `json:"usage_count"` // Number of containers using this volume
}
