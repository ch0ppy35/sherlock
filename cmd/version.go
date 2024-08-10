package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display sherlock's version",
	Long:  "The version command will display sherlock's version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s  Commit: %s  Date: %s", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
