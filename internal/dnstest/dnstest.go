package dnstest

import (
	"fmt"
	"sync"

	cfg "github.com/ch0ppy35/sherlock/internal/config"
	"github.com/ch0ppy35/sherlock/internal/dns"
)

// RunAllTestsInConfig checks if the DNS records for the hosts in the configuration
// match the expected values specified in the config and reports all errors at the end.
func RunAllTestsInConfig(config cfg.Config, client dns.TinyDNSClient) error {
	var allErrors []error
	var mu sync.Mutex
	var wg sync.WaitGroup

	hostTests := make(map[string][]cfg.TestConfig)
	results := make(map[string]*dns.DNSRecords)
	errors := make(map[string]error)

	fmt.Printf("Using the following DNS server: %s\n\n", config.DNSServer)
	for _, test := range config.Tests {
		hostTests[test.Host] = append(hostTests[test.Host], test)
	}

	for host := range hostTests {
		wg.Add(1)
		go queryDNSForHost(host, config.DNSServer, client, results, errors, &mu, &wg)
	}
	wg.Wait()

	for host, tests := range hostTests {
		runTestsForHost(host, tests, results, errors, &allErrors)
	}

	if len(allErrors) > 0 {
		return fmt.Errorf("test failures:\n%v", allErrors)
	}
	return nil
}

func queryDNSForHost(host string, server string, client dns.TinyDNSClient, results map[string]*dns.DNSRecords, errors map[string]error, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	records, err := dns.QueryDNS(host, server, client)
	mu.Lock()
	results[host] = records
	errors[host] = err
	mu.Unlock()
}

func runTestsForHost(host string, tests []cfg.TestConfig, results map[string]*dns.DNSRecords, errors map[string]error, allErrors *[]error) {
	fmt.Printf("####################################################\n")
	fmt.Printf("Running tests for: %s...\n", host)

	if err, found := errors[host]; found && err != nil {
		fmt.Printf("Failed to query DNS for host %s: %v\n", host, err)
		*allErrors = append(*allErrors, fmt.Errorf("failed to query DNS for host %s: %w", host, err))
		return
	}

	records := results[host]

	for _, test := range tests {
		fmt.Printf("----------------------------------------------------\n")
		fmt.Printf("Testing '%s' records\n", test.TestType)
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
				actualValues = append(actualValues, fmt.Sprintf("%s %d\n", mx.Host, mx.Pref))
			}
		case "txt":
			actualValues = records.TXTRecords
		case "ns":
			actualValues = records.NSRecords
		default:
			*allErrors = append(*allErrors, fmt.Errorf("unknown test type: %s", test.TestType))
			continue
		}

		if err := dns.CompareRecords(test.ExpectedValues, actualValues); err != nil {
			fmt.Println("BAD — Records don't match the configuration")
			*allErrors = append(*allErrors, fmt.Errorf("DNS check failed for host %s: %v", host, err))
		} else {
			fmt.Printf("GOOD — All records match the configuration\n")
		}
	}
	fmt.Printf("####################################################\n\n")
}
