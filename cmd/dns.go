package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// dnsCmd represents the dns command
var dnsCmd = &cobra.Command{
	Use:                   "dns",
	DisableFlagsInUseLine: true,
	Short:                 "Sherlock DNS Testing Toolkit",
	Long: `The dns command is a toolkit designed for performing DNS-related tasks such as running 
predefined DNS tests based on a configuration file, or performing individual DNS queries and comparisons.

Examples:
  sherlock dns run --config path/to/config.yaml
  sherlock dns test --type a --host example.com --expected "10.0.0.100" --server 1.1.1.1`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(dnsCmd)
}
