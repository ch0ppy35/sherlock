package config

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		configFile  string
		expected    Config
		expectError bool
	}{
		{
			name:       "Valid Config File",
			configFile: "testdata/valid_config.yaml",
			expected: Config{
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
			configFile:  "testdata/missing_config.yaml",
			expectError: true,
		},
		{
			name:       "Missing DNS Server",
			configFile: "testdata/missing_dns_server.yaml",
			expected: Config{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := LoadConfig(tt.configFile)

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
