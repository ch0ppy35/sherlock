package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sherlock",
	Short: "A CLI tool for testing DNS records",
	Long: `sherlock is a command-line tool designed to perform DNS record tests
based on a specified configuration file. It allows you to run various types
of DNS checks, such as verifying A, AAAA, CNAME, MX, TXT, and NS records.

This tool supports a configuration file in YAML format where you can define
the expected DNS records for different hosts. The 'run' command executes all
tests defined in the configuration and provides a summary of any discrepancies
found.

Usage examples:
  sherlock run --config path/to/config.yaml

For more information on individual commands, use the --help flag with the
command name.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sherlock.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
