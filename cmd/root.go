package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// OutputJSON determines if output should be structured as JSON.
	OutputJSON bool
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "dockit",
	Short: "A safe, intelligent, audit-first Docker disk analysis CLI",
	Long: `dockit provides full visibility into your Docker disk usage
and offers risk-aware cleanup recommendations.

It replaces the blind usage of 'docker system prune -a'`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVar(&OutputJSON, "json", false, "Output results in JSON format")
}
