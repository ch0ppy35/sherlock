package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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

	// Create test files for the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.configFile == "testdata/empty_config.yaml" {
				// Create an empty file
				f, err := os.Create(tt.configFile)
				if err != nil {
					t.Fatalf("failed to create empty config file: %v", err)
				}
				f.Close()
			}

			config, err := LoadConfig(tt.configFile)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, config)
			}

			if tt.configFile == "testdata/empty_config.yaml" {
				os.Remove(tt.configFile)
			}
		})
	}
}
