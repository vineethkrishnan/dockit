package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	appVersion = "dev"
	appCommit  = "none"
	appDate    = "unknown"
)

// SetVersion is injected from main.go
func SetVersion(v, c, d string) {
	appVersion = v
	appCommit = c
	appDate = d
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of dockit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dockit cli - version %s (commit: %s, built at: %s)\n", appVersion, appCommit, appDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
