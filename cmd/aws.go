package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// awsCmd represents the aws command
var awsCmd = &cobra.Command{
	Use:                   "aws",
	DisableFlagsInUseLine: true,
	Short:                 "Sherlock AWS Toolkit",
	Long:                  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(awsCmd)
}
