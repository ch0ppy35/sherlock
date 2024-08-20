package dns_test_executor

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"net"
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
			expectedError: "test failures:\n[failed to query DNS for host example.com: failed to query DNS records: network error]",
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
			client := &dns.MockIDNSClient{
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
			executor := NewDNSTestExecutor(tt.config, client)
			err := executor.RunAllTests()

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
			client := &dns.MockIDNSClient{
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

			executor := NewDNSTestExecutor(cfg.Config{}, client)
			var wg sync.WaitGroup

			wg.Add(1)
			executor.queryDNSForHost(tt.host, &wg)
			wg.Wait()

			if !reflect.DeepEqual(executor.Results[tt.host], tt.expected) {
				t.Errorf("queryDNSForHost() records = %v, expected %v", executor.Results[tt.host], tt.expected)
			}

			if tt.expectedError == "" {
				if executor.Errors[tt.host] != nil {
					t.Errorf("queryDNSForHost() unexpected error = %v", executor.Errors[tt.host])
				}
			} else if executor.Errors[tt.host] == nil || !strings.Contains(executor.Errors[tt.host].Error(), tt.expectedError) {
				t.Errorf("queryDNSForHost() error = %v, expectedError %v", executor.Errors[tt.host], tt.expectedError)
			}
		})
	}
}

func Test_getDNSRecords(t *testing.T) {
	tests := []struct {
		name     string
		testType string
		records  *dns.DNSRecords
		expected []string
	}{
		{
			name:     "Valid A records",
			testType: "a",
			records: &dns.DNSRecords{
				ARecords: []string{"10.0.0.1", "10.0.0.2"},
			},
			expected: []string{"10.0.0.1", "10.0.0.2"},
		},
		{
			name:     "Valid AAAA records",
			testType: "aaaa",
			records: &dns.DNSRecords{
				AAAARecords: []string{"2001:db8::1", "2001:db8::2"},
			},
			expected: []string{"2001:db8::1", "2001:db8::2"},
		},
		{
			name:     "Valid CNAME records",
			testType: "cname",
			records: &dns.DNSRecords{
				CNAMERecords: []string{"alias.example.com"},
			},
			expected: []string{"alias.example.com"},
		},
		{
			name:     "Valid MX records",
			testType: "mx",
			records: &dns.DNSRecords{
				MXRecords: []dns.MXRecord{
					{Host: "mail.example.com", Pref: 10},
					{Host: "backup.example.com", Pref: 20},
				},
			},
			expected: []string{"mail.example.com 10", "backup.example.com 20"},
		},
		{
			name:     "Valid TXT records",
			testType: "txt",
			records: &dns.DNSRecords{
				TXTRecords: []string{"v=spf1 include:_spf.example.com ~all"},
			},
			expected: []string{"v=spf1 include:_spf.example.com ~all"},
		},
		{
			name:     "Valid NS records",
			testType: "ns",
			records: &dns.DNSRecords{
				NSRecords: []string{"ns1.example.com", "ns2.example.com"},
			},
			expected: []string{"ns1.example.com", "ns2.example.com"},
		},
		{
			name:     "Empty records for A type",
			testType: "a",
			records: &dns.DNSRecords{
				ARecords: []string{},
			},
			expected: []string{},
		},
		{
			name:     "Empty records for AAAA type",
			testType: "aaaa",
			records: &dns.DNSRecords{
				AAAARecords: []string{},
			},
			expected: []string{},
		},
		{
			name:     "Empty records for CNAME type",
			testType: "cname",
			records: &dns.DNSRecords{
				CNAMERecords: []string{},
			},
			expected: []string{},
		},
		{
			name:     "Empty records for MX type",
			testType: "mx",
			records: &dns.DNSRecords{
				MXRecords: []dns.MXRecord{},
			},
			expected: []string{},
		},
		{
			name:     "Empty records for TXT type",
			testType: "txt",
			records: &dns.DNSRecords{
				TXTRecords: []string{},
			},
			expected: []string{},
		},
		{
			name:     "Empty records for NS type",
			testType: "ns",
			records: &dns.DNSRecords{
				NSRecords: []string{},
			},
			expected: []string{},
		},
		{
			name:     "Unhandled test type",
			testType: "srv",
			records: &dns.DNSRecords{
				ARecords: []string{"10.0.0.1"},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := &DNSTestExecutor{}
			result := executor.getDNSRecords(tt.testType, tt.records)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("getDNSRecords() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
