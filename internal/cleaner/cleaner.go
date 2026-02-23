package cleaner

import (
	"context"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/dustin/go-humanize"
	"github.com/vineethkrishnan/dockit/internal/analyzer"
	"github.com/vineethkrishnan/dockit/pkg/models"
)

// Deleter defines the interface required to remove Docker resources.
type Deleter interface {
	RemoveContainer(ctx context.Context, id string) error
	RemoveImage(ctx context.Context, id string) error
	RemoveVolume(ctx context.Context, name string) error
}

// targetResource is a generic struct for items we want to delete to make the UI uniform.
type targetResource struct {
	ID    string
	Type  string // "Container", "Image", "Volume"
	Name  string
	Size  int64
	Score models.Score
}

// CleanupPlan organizes what needs to happen.
type CleanupPlan struct {
	Targets        []targetResource
	TotalTargets   int
	SpaceToReclaim int64
}

// BuildPlan takes correlated data and extracts anything SAFE or REVIEW.
func BuildPlan(data *analyzer.CorrelatedData) *CleanupPlan {
	plan := &CleanupPlan{}

	for _, c := range data.Containers {
		if c.Score == models.ScoreProtected {
			continue
		}
		name := c.ID[:10]
		if len(c.Names) > 0 {
			name = c.Names[0]
		}
		plan.Targets = append(plan.Targets, targetResource{
			ID:    c.ID,
			Type:  "Container",
			Name:  name,
			Size:  c.SizeRw,
			Score: c.Score,
		})
		plan.SpaceToReclaim += c.SizeRw
	}

	for _, i := range data.Images {
		if i.Score == models.ScoreProtected {
			continue
		}
		name := "<none>"
		if len(i.RepoTags) > 0 {
			name = i.RepoTags[0]
		}
		plan.Targets = append(plan.Targets, targetResource{
			ID:    i.ID,
			Type:  "Image",
			Name:  name,
			Size:  i.Size,
			Score: i.Score,
		})
		plan.SpaceToReclaim += i.Size
	}

	for _, v := range data.Volumes {
		if v.Score == models.ScoreProtected {
			continue
		}
		plan.Targets = append(plan.Targets, targetResource{
			ID:    v.Name, // Volumes use Name as ID
			Type:  "Volume",
			Name:  v.Name,
			Size:  v.Size,
			Score: v.Score,
		})
		plan.SpaceToReclaim += v.Size
	}

	plan.TotalTargets = len(plan.Targets)
	return plan
}

// PrintDryRun outputs what would have been deleted.
func (p *CleanupPlan) PrintDryRun() {
	if p.TotalTargets == 0 {
		fmt.Println("No SAFE or REVIEW resources found to clean. Disk is healthy.")
		return
	}

	fmt.Println("\n--- DRY RUN: Cleanup Plan ---")
	fmt.Printf("dockit identified %d resources that are eligible for deletion.\n", p.TotalTargets)
	fmt.Println("Only SAFE and REVIEW items are included. PROTECTED items are ignored.")
	fmt.Printf("Total space that would be reclaimed: %s\n\n", humanize.Bytes(uint64(p.SpaceToReclaim)))

	fmt.Printf("%-15s %-40s %-10s %s\n", "TYPE", "ID/NAME", "SCORE", "SIZE")
	for _, t := range p.Targets {
		name := t.Name
		if len(name) > 38 {
			name = name[:36] + ".."
		}
		fmt.Printf("%-15s %-40s %-10s %s\n", t.Type, name, t.Score, humanize.Bytes(uint64(t.Size)))
	}
	fmt.Println("\nTo apply these deletions, run: dockit clean --apply")
}

// Execute asks for confirmation and deletes the resources.
func (p *CleanupPlan) Execute(ctx context.Context, client Deleter) {
	if p.TotalTargets == 0 {
		fmt.Println("No resources to clean.")
		return
	}

	// Always require interactive confirmation.
	confirm := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("dockit will permanently delete %d resources and reclaim %s. Are you sure?", p.TotalTargets, humanize.Bytes(uint64(p.SpaceToReclaim))),
		Default: false,
	}
	err := survey.AskOne(prompt, &confirm)
	if err != nil || !confirm {
		fmt.Println("Cleanup aborted.")
		os.Exit(0)
	}

	fmt.Println("\nCleaning resources...")
	successCount := 0
	freedSpace := int64(0)

	for _, target := range p.Targets {
		var err error
		switch target.Type {
		case "Container":
			err = client.RemoveContainer(ctx, target.ID)
		case "Image":
			err = client.RemoveImage(ctx, target.ID)
		case "Volume":
			err = client.RemoveVolume(ctx, target.ID)
		}

		if err != nil {
			fmt.Printf("❌ Failed to delete %s %s: %v\n", target.Type, target.Name, err)
		} else {
			fmt.Printf("✅ Deleted %s %s (freed %s)\n", target.Type, target.Name, humanize.Bytes(uint64(target.Size)))
			successCount++
			freedSpace += target.Size
		}
	}

	fmt.Printf("\nCleanup complete! Deleted %d/%d resources. Reclaimed %s.\n", successCount, p.TotalTargets, humanize.Bytes(uint64(freedSpace)))
}
