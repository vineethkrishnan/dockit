package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"github.com/vineethkrishnan/dockit/internal/dockerclient"
	"github.com/vineethkrishnan/dockit/internal/logger"
	"github.com/vineethkrishnan/dockit/pkg/models"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Detect runaway container logs silently filling up your disk",
	Long: `Inspects all containers to find their actual json-file log size on disk.
Warns if containers are generating excessive logs without a max-size rotation policy.`,
	RunE: runLogs,
}

func init() {
	rootCmd.AddCommand(logsCmd)
}

func runLogs(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	client, err := dockerclient.NewClient()
	if err != nil {
		return fmt.Errorf("failed to connect to docker daemon: %w", err)
	}

	fmt.Println("Finding container log paths on disk...")
	rawMetrics, err := client.GetLogMetrics(ctx)
	if err != nil {
		return err
	}

	engine := logger.NewEngine(logger.DefaultConfig)
	sortedMetrics, totalSize := engine.AnalyzeLogSizes(rawMetrics)

	if OutputJSON {
		return printLogsJSON(sortedMetrics)
	}

	printLogsTable(engine, sortedMetrics, totalSize)
	return nil
}

func printLogsJSON(data []*models.LogMetrics) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func printLogsTable(engine *logger.Engine, data []*models.LogMetrics, totalSize int64) {
	fmt.Printf("\n--- CONTAINER LOG SIZES (Total: %s) ---\n", humanize.Bytes(uint64(totalSize)))
	fmt.Printf("%-20s %-15s %s\n", "CONTAINER", "SIZE", "WARNINGS")
	
	for _, m := range data {
		name := m.ContainerName
		if len(name) > 18 {
			name = name[:18] + ".."
		}

		warning := ""
		if engine.IsExcessive(m.LogSize) {
			warning = "🚨 EXCESSIVE - Consider adding 'log-opt max-size=10m'"
		} else if m.LogSize == 0 && m.LogPath == "" {
			warning = "Log path not found locally"
		}

		fmt.Printf("%-20s %-15s %s\n", name, humanize.Bytes(uint64(m.LogSize)), warning)
	}
}
