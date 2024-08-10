package cmd

import (
	"fmt"
	"os"

	cfg "github.com/ch0ppy35/sherlock/internal/config"
	"github.com/ch0ppy35/sherlock/internal/dnstest"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

var configFile string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run DNS tests based on the provided configuration",
	Long: `Run DNS tests for the hosts specified in the configuration file. The command will
execute all tests defined in the config and only report errors at the end. This ensures that
all tests are run in a single execution without needing to rerun the tool.

The configuration file should be in YAML format and specify the expected DNS records for
the hosts to be tested. The tests will verify if the actual DNS records match the expected
values and report any discrepancies after all tests are completed.

Example usage:
  sherlock run --config path/to/config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := cfg.LoadConfig(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		client := new(dns.Client)
		err = dnstest.RunAllTestsInConfig(config, client)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running tests: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (should be in config/config.yaml)")
	runCmd.MarkPersistentFlagRequired("config")
}
