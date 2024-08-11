package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "N/A"
	date    = "NOW"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display Sherlock's version",
	Long:  "The version command will display Sherlock's version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Sherlock\nVersion: %s\nCommit: %s\nBuild Time: %s\n", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
