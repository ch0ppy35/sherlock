package cmd

import (
	"fmt"
	"os"

	"github.com/ch0ppy35/sherlock/internal/dns"
	"github.com/ch0ppy35/sherlock/internal/ui"
	d "github.com/miekg/dns"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:                   "test --type <a|aaaa|cname|mx|txt|ns> --host <hostname> --expected <record1,record2,...> --server <dns-server>",
	DisableFlagsInUseLine: true,
	Example:               `sherlock test --type a --host example.com --expected "10.0.0.100" --server 1.1.1.1`,
	Short:                 "Run DNS tests based on provided parameters",
	Long: `The test command allows you to perform DNS queries for specific record types and 
compare the results with the expected values.

This command is useful for testing individual DNS records without needing to run 
the suite of tests defined in a config file. It lets you specify the 
DNS server, the type of record to query, the expected values, and the host you 
want to look up. It performs the query and then compares the actual results against 
the expected values.

Flags:
	--type string   The type of DNS record to query (e.g., a, aaaa, cname, mx, txt, ns)
	--host string   The hostname to look up (e.g., example.com)
	--expected string   Comma-separated list of expected DNS records
	--server string   The DNS server to query (e.g., 1.1.1.1)`,
	Run: func(cmd *cobra.Command, args []string) {
		testType, expectedValues, dnsServer, host, err := parseFlags(cmd)
		if err != nil {
			ui.PrintMsgWithStatus("ERROR", "red", "Error: %v\n\n", err)
			cmd.Usage()
			os.Exit(1)
		}

		if err := runDNSQueryAndCompare(testType, expectedValues, dnsServer, host); err != nil {
			ui.PrintMsgWithStatus("FAIL", "red", "Test failed: %v\n", err)
			os.Exit(1)
		} else {
			ui.PrintMsgWithStatus("GOOD", "green", "Record matches!\n")
		}
	},
}

func parseFlags(cmd *cobra.Command) (string, []string, string, string, error) {
	testType, _ := cmd.Flags().GetString("type")
	expectedValues, _ := cmd.Flags().GetStringSlice("expected")
	dnsServer, _ := cmd.Flags().GetString("server")
	host, _ := cmd.Flags().GetString("host")

	if testType == "" || len(expectedValues) == 0 || dnsServer == "" || host == "" {
		return "", nil, "", "", fmt.Errorf("all flags (--type, --host, --expected, --server) are required")
	}

	return testType, expectedValues, dnsServer, host, nil
}

func runDNSQueryAndCompare(testType string, expectedValues []string, dnsServer, domain string) error {
	client := new(d.Client)

	actualValues, err := dns.QueryAndExtract(client, testType, dnsServer, domain)
	if err != nil {
		return fmt.Errorf("error querying DNS: %v", err)
	}

	if err := dns.CompareRecords(expectedValues, actualValues); err != nil {
		return fmt.Errorf("DNS comparison failed: %v", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().StringP("type", "t", "", "DNS record type (e.g., a, aaaa, cname, mx, txt, ns)")
	testCmd.Flags().StringSliceP("expected", "e", []string{}, "Expected DNS records, comma-separated")
	testCmd.Flags().StringP("server", "s", "", "DNS server to query (e.g., 1.1.1.1)")
	testCmd.Flags().StringP("host", "H", "", "The host you want to look up (e.g., example.com)")
}
