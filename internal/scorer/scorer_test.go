package scorer

import (
	"testing"
	"time"

	"github.com/vineethkrishnan/dockit/pkg/models"
)

func TestScoreContainer(t *testing.T) {
	scorer := NewScorer(DefaultConfig)

	// Test Running
	c1 := &models.Container{State: "running"}
	scorer.ScoreContainer(c1)
	if c1.Score != models.ScoreProtected {
		t.Errorf("expected ScoreProtected for running container, got %s", c1.Score)
	}

	// Test Recently Stopped
	c2 := &models.Container{State: "exited", Created: time.Now().Add(-24 * time.Hour)}
	scorer.ScoreContainer(c2)
	if c2.Score != models.ScoreReview {
		t.Errorf("expected ScoreReview for recent exited container, got %s", c2.Score)
	}

	// Test Old Stopped
	c3 := &models.Container{State: "exited", Created: time.Now().Add(-10 * 24 * time.Hour)}
	scorer.ScoreContainer(c3)
	if c3.Score != models.ScoreSafe {
		t.Errorf("expected ScoreSafe for old exited container, got %s", c3.Score)
	}
}

func TestScoreImage(t *testing.T) {
	scorer := NewScorer(DefaultConfig)

	// In Use
	i1 := &models.Image{Containers: 1}
	scorer.ScoreImage(i1)
	if i1.Score != models.ScoreProtected {
		t.Errorf("expected ScoreProtected for in-use image, got %s", i1.Score)
	}

	// Unused but Named
	i2 := &models.Image{Containers: 0, Dangling: false}
	scorer.ScoreImage(i2)
	if i2.Score != models.ScoreReview {
		t.Errorf("expected ScoreReview for unused named image, got %s", i2.Score)
	}

	// Dangling and Old
	i3 := &models.Image{Containers: 0, Dangling: true, Created: time.Now().Add(-10 * 24 * time.Hour)}
	scorer.ScoreImage(i3)
	if i3.Score != models.ScoreSafe {
		t.Errorf("expected ScoreSafe for old dangling image, got %s", i3.Score)
	}
}
