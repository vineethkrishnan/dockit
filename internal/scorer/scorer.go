package scorer

import (
	"time"

	"github.com/vineethkrishnan/dockit/pkg/models"
)

// Config allows varying the time thresholds for heuristic scoring.
type Config struct {
	ReviewDays int // Containers/Images younger than this are marked for REVIEW instead of SAFE.
}

// DefaultConfig provides standard conservative threshholds.
var DefaultConfig = Config{
	ReviewDays: 7,
}

// Scorer assigns risk classifications.
type Scorer struct {
	Config Config
}

func NewScorer(cfg Config) *Scorer {
	return &Scorer{Config: cfg}
}

// ScoreContainer assigns a risk class to a container.
func (s *Scorer) ScoreContainer(c *models.Container) {
	if c.State == "running" || c.State == "restarting" || c.State == "paused" {
		c.Score = models.ScoreProtected
		c.Reason = "Container is currently active"
		return
	}

	age := time.Since(c.Created)
	threshold := time.Duration(s.Config.ReviewDays) * 24 * time.Hour

	if age < threshold {
		c.Score = models.ScoreReview
		c.Reason = "Container stopped, but created recently"
		return
	}

	c.Score = models.ScoreSafe
	c.Reason = "Container is stopped and old"
}

// ScoreImage assigns a risk class to an image.
func (s *Scorer) ScoreImage(img *models.Image) {
	if img.Containers > 0 {
		img.Score = models.ScoreProtected
		img.Reason = "Image is currently backing a container"
		return
	}

	age := time.Since(img.Created)
	threshold := time.Duration(s.Config.ReviewDays) * 24 * time.Hour

	if !img.Dangling {
		// Named image, but no active containers
		img.Score = models.ScoreReview
		img.Reason = "Image specifies a repository/tag but is unused"
		return
	}

	// Dangling images
	if age < threshold {
		img.Score = models.ScoreReview
		img.Reason = "Image is dangling but was created recently"
		return
	}

	img.Score = models.ScoreSafe
	img.Reason = "Image is dangling and old"
}

// ScoreVolume assigns a risk class to a volume.
func (s *Scorer) ScoreVolume(vol *models.Volume) {
	if vol.UsageCount > 0 {
		vol.Score = models.ScoreProtected
		vol.Reason = "Volume is attached to a container"
		return
	}

	// Volumes that have 0 usage but might not have a reliable CreatedAt
	// For MVP, if it isn't attached to anything, it's SAFE to delete.
	// Users can skip if they want to retain data.
	vol.Score = models.ScoreSafe
	vol.Reason = "Volume is orphaned (not attached to any container)"
}
