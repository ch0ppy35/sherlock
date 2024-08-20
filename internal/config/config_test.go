package config

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		configFile  string
		expected    DNSRecordsFullTestConfig
		expectError bool
	}{
		{
			name:       "Valid Config File",
			configFile: "dnstestdata/valid_config.yaml",
			expected: DNSRecordsFullTestConfig{
				DNSServer: "8.8.8.8",
				Tests: []DNSTestConfig{
					{
						ExpectedValues: []string{"1.1.1.1"},
						Host:           "example.com",
						TestType:       "A",
					},
					{
						ExpectedValues: []string{"text value"},
						Host:           "example.com",
						TestType:       "TXT",
					},
				},
			},
			expectError: false,
		},
		{
			name:        "Empty Config File Path",
			configFile:  "",
			expectError: true,
		},
		{
			name:        "Missing Config File",
			configFile:  "dnstestdata/missing_config.yaml",
			expectError: true,
		},
		{
			name:       "Missing DNS Server",
			configFile: "dnstestdata/missing_dns_server.yaml",
			expected: DNSRecordsFullTestConfig{
				DNSServer: "1.1.1.1", // Default DNS Server
				Tests: []DNSTestConfig{
					{
						ExpectedValues: []string{"1.1.1.1"},
						Host:           "example.com",
						TestType:       "A",
					},
				},
			},
			expectError: false,
		},
		{
			name:        "No Tests Defined",
			configFile:  "dnstestdata/no_tests_defined.yaml",
			expectError: true,
		},
		{
			name:        "Missing Expected Values",
			configFile:  "dnstestdata/missing_expected_values.yaml",
			expectError: true,
		},
		{
			name:        "Missing Host",
			configFile:  "dnstestdata/missing_host.yaml",
			expectError: true,
		},
		{
			name:        "Missing Test Type",
			configFile:  "dnstestdata/missing_test_type.yaml",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := LoadDNSRecordsFullTestConfig(tt.configFile)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("didn't expect an error but got %v", err)
				return
			}

			if !reflect.DeepEqual(tt.expected, config) {
				t.Errorf("expected config %+v, but got %+v", tt.expected, config)
			}
		})
	}
}
