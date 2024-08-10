package dnstest

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"net"
	"reflect"
	"strings"

	cfg "github.com/ch0ppy35/sherlock/internal/config"
	"github.com/ch0ppy35/sherlock/internal/dns"
	d "github.com/miekg/dns"
)

func Test_RunAllTestsInConfig(t *testing.T) {
	tests := []struct {
		name          string
		config        cfg.Config
		mockResponses map[uint16]*d.Msg
		mockError     error
		expectedError string
	}{
		{
			name: "Valid configuration",
			config: cfg.Config{
				DNSServer: "8.8.8.8",
				Tests: []cfg.DNSTestConfig{
					{
						Host:           "example.com",
						TestType:       "a",
						ExpectedValues: []string{"10.0.0.1"},
					},
				},
			},
			mockResponses: map[uint16]*d.Msg{
				d.TypeA: {
					Answer: []d.RR{
						&d.A{Hdr: d.RR_Header{Name: "example.com."}, A: net.ParseIP("10.0.0.1")},
					},
				},
			},
			expectedError: "",
		},
		{
			name: "Configuration with missing records",
			config: cfg.Config{
				DNSServer: "8.8.8.8",
				Tests: []cfg.DNSTestConfig{
					{
						Host:           "example.com",
						TestType:       "a",
						ExpectedValues: []string{"10.0.0.2"},
					},
				},
			},
			mockResponses: map[uint16]*d.Msg{
				d.TypeA: {
					Answer: []d.RR{
						&d.A{Hdr: d.RR_Header{Name: "example.com."}, A: net.ParseIP("10.0.0.1")},
					},
				},
			},
			expectedError: "test failures:\n[DNS check failed for host example.com: mismatched records found]",
		},
		{
			name: "Configuration with query error",
			config: cfg.Config{
				DNSServer: "8.8.8.8",
				Tests: []cfg.DNSTestConfig{
					{
						Host:           "example.com",
						TestType:       "a",
						ExpectedValues: []string{"10.0.0.1"},
					},
				},
			},
			mockError:     fmt.Errorf("network error"),
			expectedError: "test failures:\n[failed to query DNS for host example.com: failed to query DNS records for type 1: network error]",
		},
		{
			name: "Empty configuration",
			config: cfg.Config{
				DNSServer: "8.8.8.8",
				Tests:     []cfg.DNSTestConfig{},
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &dns.MockTinyDNSClient{
				MockExchange: func(msg *d.Msg, server string) (*d.Msg, time.Duration, error) {
					if tt.mockError != nil {
						return nil, 0, tt.mockError
					}
					if resp, ok := tt.mockResponses[msg.Question[0].Qtype]; ok {
						return resp, 0, nil
					}
					return &d.Msg{}, 0, nil
				},
			}

			err := RunAllTestsInConfig(tt.config, client)

			if (err != nil && tt.expectedError == "") || (err == nil && tt.expectedError != "") {
				t.Errorf("RunAllTestsInConfig() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			if err != nil && !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("RunAllTestsInConfig() error = '%v', expectedError '%v'", err, tt.expectedError)
			}
		})
	}
}

func Test_queryDNSForHost(t *testing.T) {
	tests := []struct {
		name          string
		host          string
		server        string
		mockResponses map[uint16]*d.Msg
		mockError     error
		expected      *dns.DNSRecords
		expectedError string
	}{
		{
			name:   "Valid A record query",
			host:   "example.com",
			server: "8.8.8.8",
			mockResponses: map[uint16]*d.Msg{
				d.TypeA: {
					Answer: []d.RR{
						&d.A{Hdr: d.RR_Header{Name: "example.com."}, A: net.ParseIP("10.0.0.1")},
					},
				},
			},
			expected: &dns.DNSRecords{
				ARecords: []string{"10.0.0.1"},
			},
		},
		{
			name:   "Query returns no answers",
			host:   "example.com",
			server: "8.8.8.8",
			mockResponses: map[uint16]*d.Msg{
				d.TypeA: {},
			},
			expected: &dns.DNSRecords{},
		},
		{
			name:          "Query returns an error",
			host:          "example.com",
			server:        "8.8.8.8",
			mockError:     fmt.Errorf("network error"),
			expectedError: "network error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &dns.MockTinyDNSClient{
				MockExchange: func(msg *d.Msg, server string) (*d.Msg, time.Duration, error) {
					if tt.mockError != nil {
						return nil, 0, tt.mockError
					}
					if resp, ok := tt.mockResponses[msg.Question[0].Qtype]; ok {
						return resp, 0, nil
					}
					return &d.Msg{}, 0, nil
				},
			}

			results := make(map[string]*dns.DNSRecords)
			errors := make(map[string]error)
			var mu sync.Mutex
			var wg sync.WaitGroup

			wg.Add(1)
			queryDNSForHost(tt.host, tt.server, client, results, errors, &mu, &wg)
			wg.Wait()

			mu.Lock()
			defer mu.Unlock()

			if !reflect.DeepEqual(results[tt.host], tt.expected) {
				t.Errorf("queryDNSForHost() records = %v, expected %v", results[tt.host], tt.expected)
			}

			if tt.expectedError == "" {
				if errors[tt.host] != nil {
					t.Errorf("queryDNSForHost() unexpected error = %v", errors[tt.host])
				}
			} else if errors[tt.host] == nil || !strings.Contains(errors[tt.host].Error(), tt.expectedError) {
				t.Errorf("queryDNSForHost() error = %v, expectedError %v", errors[tt.host], tt.expectedError)
			}
		})
	}
}

func Test_runTestsForHost(t *testing.T) {
	tests := []struct {
		name           string
		host           string
		DNSTestConfigs []cfg.DNSTestConfig
		results        map[string]*dns.DNSRecords
		errors         map[string]error
		expectedErrors []error
	}{
		{
			name: "Valid AAAA record match",
			host: "example.com",
			DNSTestConfigs: []cfg.DNSTestConfig{
				{
					TestType:       "a",
					ExpectedValues: []string{"10.0.0.1"},
				},
			},
			results: map[string]*dns.DNSRecords{
				"example.com": {
					ARecords: []string{"10.0.0.1"},
				},
			},
		},
		{
			name: "Valid AAAA record match",
			host: "example.com",
			DNSTestConfigs: []cfg.DNSTestConfig{
				{
					TestType:       "aaaa",
					ExpectedValues: []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
				},
			},
			results: map[string]*dns.DNSRecords{
				"example.com": {
					AAAARecords: []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
				},
			},
			errors:         map[string]error{},
			expectedErrors: nil,
		},
		{
			name: "Valid CNAME record match",
			host: "example.com",
			DNSTestConfigs: []cfg.DNSTestConfig{
				{
					TestType:       "cname",
					ExpectedValues: []string{"cname.example.com."},
				},
			},
			results: map[string]*dns.DNSRecords{
				"example.com": {
					CNAMERecords: []string{"cname.example.com."},
				},
			},
			errors:         map[string]error{},
			expectedErrors: nil,
		},
		{
			name: "Valid MX record match",
			host: "example.com",
			DNSTestConfigs: []cfg.DNSTestConfig{
				{
					TestType:       "mx",
					ExpectedValues: []string{"mail.example.com. 10\n"},
				},
			},
			results: map[string]*dns.DNSRecords{
				"example.com": {
					MXRecords: []dns.MXRecord{
						{Host: "mail.example.com.", Pref: 10},
					},
				},
			},
			errors:         map[string]error{},
			expectedErrors: nil,
		},
		{
			name: "Valid TXT record match",
			host: "example.com",
			DNSTestConfigs: []cfg.DNSTestConfig{
				{
					TestType:       "txt",
					ExpectedValues: []string{"v=spf1 include:_spf.example.com ~all"},
				},
			},
			results: map[string]*dns.DNSRecords{
				"example.com": {
					TXTRecords: []string{"v=spf1 include:_spf.example.com ~all"},
				},
			},
			errors:         map[string]error{},
			expectedErrors: nil,
		},
		{
			name: "Valid NS record match",
			host: "example.com",
			DNSTestConfigs: []cfg.DNSTestConfig{
				{
					TestType:       "ns",
					ExpectedValues: []string{"ns1.example.com."},
				},
			},
			results: map[string]*dns.DNSRecords{
				"example.com": {
					NSRecords: []string{"ns1.example.com."},
				},
			},
			errors:         map[string]error{},
			expectedErrors: nil,
		},
		{
			name: "Records don't match configuration",
			host: "example.com",
			DNSTestConfigs: []cfg.DNSTestConfig{
				{
					TestType:       "a",
					ExpectedValues: []string{"10.0.0.2"},
				},
			},
			results: map[string]*dns.DNSRecords{
				"example.com": {
					ARecords: []string{"10.0.0.1"},
				},
			},
			errors: map[string]error{},
			expectedErrors: []error{
				fmt.Errorf("DNS check failed for host example.com: mismatched records found"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allErrors := []error{}
			runTestsForHost(tt.host, tt.DNSTestConfigs, tt.results, tt.errors, &allErrors)

			if len(allErrors) != len(tt.expectedErrors) {
				t.Errorf("runTestsForHost() allErrors = %v, expectedErrors %v", allErrors, tt.expectedErrors)
				return
			}

			for i, err := range allErrors {
				if err.Error() != tt.expectedErrors[i].Error() {
					t.Errorf("runTestsForHost() error = '%v', expectedError '%v'", err, tt.expectedErrors[i])
				}
			}
		})
	}
}
