package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"github.com/vineethkrishnan/dockit/internal/dockerclient"
	"github.com/vineethkrishnan/dockit/pkg/models"
)

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Provides a high-level overview of Docker disk usage",
	Long: `Fetches the total disk space consumed by Docker 
Images, Containers, Volumes, and Build Cache.`,
	RunE: runSummary,
}

func init() {
	rootCmd.AddCommand(summaryCmd)
}

func runSummary(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := dockerclient.NewClient()
	if err != nil {
		return fmt.Errorf("failed to connect to docker daemon: %w", err)
	}

	if err := client.Ping(ctx); err != nil {
		return fmt.Errorf("docker daemon is not responsive: %w", err)
	}

	summary, err := client.GetDiskSummary(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch disk summary: %w", err)
	}

	if OutputJSON {
		return printJSON(summary)
	}

	printTable(summary)
	return nil
}

func printJSON(summary *models.DiskSummary) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(summary)
}

func printTable(s *models.DiskSummary) {
	fmt.Println("TYPE\t\tTOTAL\t\tACTIVE\t\tSIZE\t\tRECLAIMABLE")
	fmt.Printf("Images\t\t%d\t\t%d\t\t%s\t\t%s\n",
		s.Images.TotalCount, s.Images.ActiveCount, humanize.Bytes(uint64(s.Images.TotalSize)), humanize.Bytes(uint64(s.Images.Reclaimable)))

	fmt.Printf("Containers\t%d\t\t%d\t\t%s\t\t%s\n",
		s.Containers.TotalCount, s.Containers.ActiveCount, humanize.Bytes(uint64(s.Containers.TotalSize)), humanize.Bytes(uint64(s.Containers.Reclaimable)))

	fmt.Printf("Local Volumes\t%d\t\t%d\t\t%s\t\t%s\n",
		s.Volumes.TotalCount, s.Volumes.ActiveCount, humanize.Bytes(uint64(s.Volumes.TotalSize)), humanize.Bytes(uint64(s.Volumes.Reclaimable)))

	fmt.Printf("Build Cache\t%d\t\t%d\t\t%s\t\t%s\n",
		s.BuildCache.TotalCount, s.BuildCache.ActiveCount, humanize.Bytes(uint64(s.BuildCache.TotalSize)), humanize.Bytes(uint64(s.BuildCache.Reclaimable)))

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("Total Space:\t\t\t\t\t%s\n", humanize.Bytes(uint64(s.TotalSize)))
	fmt.Printf("Reclaimable Space:\t\t\t\t%s\n", humanize.Bytes(uint64(s.Reclaimable)))
}
