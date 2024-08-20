package config

import (
	"fmt"

	"github.com/ch0ppy35/sherlock/internal/ui"
	"github.com/spf13/viper"
)

type DNSRecordsFullTestConfig struct {
	DNSServer string          `yaml:"dnsServer"` // Optional
	Tests     []DNSTestConfig `yaml:"tests"`     // Required
}

type DNSTestConfig struct {
	ExpectedValues []string `yaml:"expectedValues"` // Required
	Host           string   `yaml:"host"`           // Required
	TestType       string   `yaml:"testType"`       // Required
}

func (c *DNSRecordsFullTestConfig) validate() error {
	if len(c.Tests) == 0 {
		return fmt.Errorf("no tests defined in the configuration")
	}

	if c.DNSServer == "" {
		ui.PrintMsgWithStatus("WARN", "hiYellow", "DNS server not set, using Cloudflare as default\n")
		c.DNSServer = "1.1.1.1"
	}

	for i, test := range c.Tests {
		if len(test.ExpectedValues) == 0 {
			return fmt.Errorf("test %d 'expectedValues' must be set and contain at least one value", i+1)
		}
		if test.Host == "" {
			return fmt.Errorf("test %d 'host' must be set", i+1)
		}
		if test.TestType == "" {
			return fmt.Errorf("test %d 'testType' must be set", i+1)
		}
	}
	return nil
}

func LoadDNSRecordsFullTestConfig(configFile string) (DNSRecordsFullTestConfig, error) {
	var config DNSRecordsFullTestConfig
	if configFile == "" {
		return DNSRecordsFullTestConfig{}, fmt.Errorf("no config file specified")
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return DNSRecordsFullTestConfig{}, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return DNSRecordsFullTestConfig{}, fmt.Errorf("unable to decode into struct: %w", err)
	}

	if err := config.validate(); err != nil {
		return DNSRecordsFullTestConfig{}, fmt.Errorf("validation issue: %w", err)
	}

	return config, nil
}
