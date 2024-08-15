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
based on a specified configuration file or individual params. It allows you to run various types
of DNS checks, such as verifying A, AAAA, CNAME, MX, TXT, and NS records.

This tool supports a configuration file in YAML format where you can define
the expected DNS records for different hosts. The 'run' command executes all
tests defined in the configuration and provides a summary of any discrepancies
found. Alternatively, you can use the 'test' command to query a DNS server for a specific
record type and compare the results with expected values directly from the command line.

Usage examples:
  sherlock run --config path/to/config.yaml
  sherlock test --type a --host example.com --expected "10.0.0.100" --server 1.1.1.1

For more information on individual commands, use the --help flag with the
command name.`,
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
}
