package cmd

import (
	"os"

	cfg "github.com/ch0ppy35/sherlock/internal/config"
	dtexc "github.com/ch0ppy35/sherlock/internal/dns_test_executor"
	"github.com/ch0ppy35/sherlock/internal/ui"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

var configFile string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:                   "run --config <path/to/config.yaml>",
	DisableFlagsInUseLine: true,
	Example:               "sherlock dns run --config path/to/config.yaml",
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
	config, err := cfg.LoadDNSRecordsFullTestConfig(configFile)
	if err != nil {
		ui.PrintMsgWithStatus("ERROR", "hiRed", "Trouble loading the config: %v\n", err)
		os.Exit(1)
	}
	client := new(dns.Client)
	executor := dtexc.NewDNSTestExecutor(config, client)
	err = executor.RunAllTests()
	if err != nil {
		ui.PrintMsgWithStatus("FAIL", "hiRed", "One or more tests failed, check above\n")
		os.Exit(1)
	}
}

func init() {
	dnsCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to the config file (config/config.yaml)")
	runCmd.MarkPersistentFlagRequired("config")
}
