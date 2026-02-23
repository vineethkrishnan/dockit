package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"github.com/vineethkrishnan/dockit/internal/analyzer"
	"github.com/vineethkrishnan/dockit/internal/dockerclient"
	"github.com/vineethkrishnan/dockit/internal/scorer"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Deep analysis of individual Docker resources",
	Long: `Fetches individual containers, images, and volumes, 
correlates them, and assigns a risk score (SAFE, REVIEW, PROTECTED)
to each resource.`,
	RunE: runAnalyze,
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := dockerclient.NewClient()
	if err != nil {
		return fmt.Errorf("failed to connect to docker daemon: %w", err)
	}

	fmt.Println("Fetching and correlating resources...")
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

	// Correlate Data
	correlated := analyzer.Analyze(containers, images, volumes)

	// Score Data
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

	if OutputJSON {
		return printAnalyzeJSON(correlated)
	}

	printAnalyzeTable(correlated)
	return nil
}

func printAnalyzeJSON(data *analyzer.CorrelatedData) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func printAnalyzeTable(data *analyzer.CorrelatedData) {
	fmt.Println("\n--- CONTAINERS ---")
	fmt.Printf("%-20s %-15s %-15s %-10s %s\n", "ID/NAME", "STATE", "SCORE", "SIZE", "REASON")
	for _, c := range data.Containers {
		name := c.ID[:10]
		if len(c.Names) > 0 {
			name = c.Names[0]
		}
		// Truncate name natively
		if len(name) > 18 {
			name = name[:18] + ".."
		}
		fmt.Printf("%-20s %-15s %-15s %-10s %s\n", name, c.State, c.Score, humanize.Bytes(uint64(c.SizeRw)), c.Reason)
	}

	fmt.Println("\n--- IMAGES ---")
	fmt.Printf("%-20s %-10s %-15s %-10s %s\n", "REPO:TAG", "DANGLING", "SCORE", "SIZE", "REASON")
	for _, img := range data.Images {
		tag := "<none>"
		if len(img.RepoTags) > 0 {
			tag = img.RepoTags[0]
		}
		if len(tag) > 18 {
			tag = tag[:18] + ".."
		}
		fmt.Printf("%-20s %-10t %-15s %-10s %s\n", tag, img.Dangling, img.Score, humanize.Bytes(uint64(img.Size)), img.Reason)
	}

	fmt.Println("\n--- VOLUMES ---")
	fmt.Printf("%-20s %-10s %-15s %-10s %s\n", "NAME", "USAGE", "SCORE", "SIZE", "REASON")
	for _, vol := range data.Volumes {
		name := vol.Name
		if len(name) > 18 {
			name = name[:18] + ".."
		}
		fmt.Printf("%-20s %-10d %-15s %-10s %s\n", name, vol.UsageCount, vol.Score, humanize.Bytes(uint64(vol.Size)), vol.Reason)
	}
}
