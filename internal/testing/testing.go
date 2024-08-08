package testing

import (
	"fmt"

	"github.com/ch0ppy35/dnsTest/internal/dns"
)

// RunAllTestsInConfig checks if the DNS records for the hosts in the configuration
// match the expected values specified in the config and reports all errors at the end.
func RunAllTestsInConfig(config Config) error {
	var allErrors []error
	hostTests := make(map[string][]TestConfig)

	// Group tests by host
	for _, test := range config.Tests {
		hostTests[test.Host] = append(hostTests[test.Host], test)
	}

	fmt.Printf("Using the following DNS server: %s\n\n", config.DNSServer)

	for host, tests := range hostTests {
		fmt.Printf("####################################################\n")
		fmt.Printf("Running tests for: %s...\n", host)

		// Query DNS records for the given host name
		records, err := dns.QueryDNS(host, config.DNSServer)
		if err != nil {
			allErrors = append(allErrors, fmt.Errorf("failed to query DNS for host %s: %w", host, err))
			continue
		}

		// Process each test for the current host
		for _, test := range tests {
			fmt.Printf("----------------------------------------------------\n")
			var actualValues []string
			switch test.TestType {
			case "a":
				actualValues = records.ARecords
			case "aaaa":
				actualValues = records.AAAARecords
			case "cname":
				actualValues = records.CNAMERecords
			case "mx":
				for _, mx := range records.MXRecords {
					actualValues = append(actualValues, fmt.Sprintf("%s %d", mx.Host, mx.Pref))
				}
			case "txt":
				actualValues = records.TXTRecords
			case "ns":
				actualValues = records.NSRecords
			default:
				allErrors = append(allErrors, fmt.Errorf("unknown test type: %s", test.TestType))
				continue
			}

			if err := dns.CompareRecords(test.ExpectedValues, actualValues); err != nil {
				allErrors = append(allErrors, fmt.Errorf("DNS check failed for host %s: %v", host, err))
			}
		}
	}
	fmt.Printf("####################################################\n\n")

	if len(allErrors) > 0 {
		return fmt.Errorf("test failures:\n%v", allErrors)
	}

	return nil
}
