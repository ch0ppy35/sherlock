package dns_test_executor

import (
	"fmt"
	"strings"
	"sync"

	cfg "github.com/ch0ppy35/sherlock/internal/config"
	"github.com/ch0ppy35/sherlock/internal/dns"
	"github.com/ch0ppy35/sherlock/internal/ui"
)

type DNSTestExecutor struct {
	Config    cfg.DNSRecordsFullTestConfig
	Client    dns.IDNSClient
	Results   map[string]*dns.DNSRecords
	Errors    map[string]error
	AllErrors []error
	mu        sync.Mutex
}

func NewDNSTestExecutor(config cfg.DNSRecordsFullTestConfig, client dns.IDNSClient) *DNSTestExecutor {
	return &DNSTestExecutor{
		Config:  config,
		Client:  client,
		Results: make(map[string]*dns.DNSRecords),
		Errors:  make(map[string]error),
	}
}

// RunAllTests executes all DNS tests defined in the configuration.
func (e *DNSTestExecutor) RunAllTests() error {
	var wg sync.WaitGroup

	hostTests := e.groupTestsByHost()

	ui.PrintMsgWithStatus("INFO", "magenta", "Using DNS server: %s\n", e.Config.DNSServer)
	for host := range hostTests {
		wg.Add(1)
		go e.queryDNSForHost(host, &wg)
	}
	wg.Wait()

	for host, tests := range hostTests {
		e.runTestsForHost(host, tests)
	}

	if len(e.AllErrors) > 0 {
		fmt.Printf("\n")
		return fmt.Errorf("test failures:\n%v", e.AllErrors)
	}
	fmt.Printf("\n")
	return nil
}

// groupTestsByHost groups DNS tests by the host.
func (e *DNSTestExecutor) groupTestsByHost() map[string][]cfg.DNSTestConfig {
	hostTests := make(map[string][]cfg.DNSTestConfig)
	for _, test := range e.Config.Tests {
		hostTests[test.Host] = append(hostTests[test.Host], test)
	}
	return hostTests
}

// queryDNSForHost queries the DNS for a specific host and stores the result.
func (e *DNSTestExecutor) queryDNSForHost(host string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer e.mu.Unlock()

	records, err := dns.QueryDNS(host, e.Config.DNSServer, e.Client)
	e.mu.Lock()

	e.Results[host] = records
	e.Errors[host] = err
}

// runTestsForHost runs all tests for a specific host.
func (e *DNSTestExecutor) runTestsForHost(host string, tests []cfg.DNSTestConfig) {
	fmt.Printf("\nRunning tests for host: %s...\n", host)

	if err, found := e.Errors[host]; found && err != nil {
		fmt.Printf("Failed to query DNS for host %s: %v\n", host, err)
		e.AllErrors = append(e.AllErrors, fmt.Errorf("failed to query DNS for host %s: %w", host, err))
		return
	}

	records := e.Results[host]
	for _, test := range tests {
		ui.PrintDashes()
		fmt.Printf("Testing '%s' records\n", test.TestType)
		actualValues := e.getDNSRecords(test.TestType, records)
		if actualValues == nil {
			fmt.Printf("Unknown test type encountered: %s for host: %s\n", test.TestType, host)
			e.AllErrors = append(e.AllErrors, fmt.Errorf("unknown test type: %s", test.TestType))
			continue
		}

		if len(actualValues) == 0 {
			fmt.Printf("No records found for test type: %s on host: %s\n", test.TestType, host)
		}

		if err := dns.CompareRecords(test.ExpectedValues, actualValues); err != nil {
			ui.PrintMsgWithStatus("BAD", "red", "Records don't match the configuration\n")
			e.AllErrors = append(e.AllErrors, fmt.Errorf("DNS check failed for host %s: %v", host, err))
		} else {
			ui.PrintMsgWithStatus("GOOD", "green", "All records match the configuration\n")
		}
	}
}

// getDNSRecords returns the relevant DNS records based on the test type.
func (e *DNSTestExecutor) getDNSRecords(testType string, records *dns.DNSRecords) []string {
	switch strings.ToLower(testType) {
	case "a":
		if len(records.ARecords) == 0 {
			return []string{}
		}
		return records.ARecords
	case "aaaa":
		if len(records.AAAARecords) == 0 {
			return []string{}
		}
		return records.AAAARecords
	case "cname":
		if len(records.CNAMERecords) == 0 {
			return []string{}
		}
		return records.CNAMERecords
	case "mx":
		var mxRecords []string
		for _, mx := range records.MXRecords {
			mxRecords = append(mxRecords, fmt.Sprintf("%s %d", mx.Host, mx.Pref))
		}
		if len(mxRecords) == 0 {
			return []string{}
		}
		return mxRecords
	case "txt":
		if len(records.TXTRecords) == 0 {
			return []string{}
		}
		return records.TXTRecords
	case "ns":
		if len(records.NSRecords) == 0 {
			return []string{}
		}
		return records.NSRecords
	default:
		fmt.Printf("Unhandled test type: %s\n", testType)
		return nil
	}
}
