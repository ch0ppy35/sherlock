package cmd

import (
	"fmt"
	"os"

	cfg "github.com/ch0ppy35/sherlock/internal/config"
	"github.com/ch0ppy35/sherlock/internal/test_executor"
	"github.com/ch0ppy35/sherlock/internal/ui"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

var configFile string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:                   "run --config <path/to/config.yaml>",
	DisableFlagsInUseLine: true,
	Example:               "sherlock run --config path/to/config.yaml",
	Short:                 "Run DNS tests based on the provided configuration",
	Long: `Run DNS tests for the hosts specified in the configuration file. The command will
execute all tests defined in the config and only report errors at the end. This ensures that
all tests are run in a single execution without needing to rerun the tool.

The configuration file should be in YAML format and specify the expected DNS records for
the hosts to be tested. The tests will verify if the actual DNS records match the expected
values and report any discrepancies after all tests are completed.`,
	Run: func(cmd *cobra.Command, args []string) {
		runTests()
	},
}

func runTests() {
	config, err := cfg.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
	client := new(dns.Client)
	executor := test_executor.NewDNSTestExecutor(config, client)
	err = executor.RunAllTests()
	if err != nil {
		ui.PrintMsgWithStatus("FAIL", "red", "One or more tests failed, check above\n")
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file (config/config.yaml)")
	runCmd.MarkPersistentFlagRequired("config")
}
