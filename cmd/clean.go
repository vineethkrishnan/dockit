package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vineethkrishnan/dockit/internal/analyzer"
	"github.com/vineethkrishnan/dockit/internal/cleaner"
	"github.com/vineethkrishnan/dockit/internal/dockerclient"
	"github.com/vineethkrishnan/dockit/internal/scorer"
)

var applyCleanup bool

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Interactively delete safe resources to free disk space",
	Long: `Cleans up Docker resources securely. 

By default, this command performs a Dry-Run to show you what would be deleted.
To actually perform the deletion, use the --apply flag.

dockit will ONLY ever delete resources marked as SAFE or REVIEW.
Running containers and attached volumes are NEVER deleted.`,
	RunE: runClean,
}

func init() {
	cleanCmd.Flags().BoolVar(&applyCleanup, "apply", false, "Actually apply deletions rather than a dry-run")
	rootCmd.AddCommand(cleanCmd)
}

func runClean(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := dockerclient.NewClient()
	if err != nil {
		return fmt.Errorf("failed to connect to docker daemon: %w", err)
	}

	fmt.Println("Analyzing disk usage...")
	containers, err := client.GetContainers(ctx)
	if err != nil {
		return err
	}
	images, err := client.GetImages(ctx)
	if err != nil {
		return err
	}
	volumes, err := client.GetVolumes(ctx)
	if err != nil {
		return err
	}

	// Correlate & Score
	correlated := analyzer.Analyze(containers, images, volumes)
	engine := scorer.NewScorer(scorer.DefaultConfig)
	for _, c := range correlated.Containers {
		engine.ScoreContainer(c)
	}
	for _, i := range correlated.Images {
		engine.ScoreImage(i)
	}
	for _, v := range correlated.Volumes {
		engine.ScoreVolume(v)
	}

	// Build the deletion plan securely
	plan := cleaner.BuildPlan(correlated)

	if !applyCleanup {
		plan.PrintDryRun()
		return nil
	}

	// Execution
	plan.Execute(ctx, client)
	return nil
}
