package config

import (
	"fmt"

	"github.com/ch0ppy35/sherlock/internal/ui"
	"github.com/spf13/viper"
)

type Config struct {
	DNSServer string          `yaml:"dnsServer"` // Optional
	Tests     []DNSTestConfig `yaml:"tests"`     // Required
}

type DNSTestConfig struct {
	ExpectedValues []string `yaml:"expectedValues"` // Required
	Host           string   `yaml:"host"`           // Required
	TestType       string   `yaml:"testType"`       // Required
}

func LoadConfig(configFile string) (Config, error) {
	var config Config
	if configFile == "" {
		return Config{}, fmt.Errorf("no config file specified")
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("unable to decode into struct: %w", err)
	}

	// Set default DNS server if not provided
	if config.DNSServer == "" {
		ui.PrintMsgWithStatus("WARN", "hiYellow", "DNS server not set, using Cloudflare as default\n")
		config.DNSServer = "1.1.1.1"
	}

	// Validate that at least one test is defined
	if len(config.Tests) == 0 {
		return Config{}, fmt.Errorf("no tests defined in the configuration")
	}

	// Validate each test to ensure all required fields are set
	for i, test := range config.Tests {
		if len(test.ExpectedValues) == 0 {
			return Config{}, fmt.Errorf("test %d 'expectedValues' must be set and contain at least one value", i+1)
		}
		if test.Host == "" {
			return Config{}, fmt.Errorf("test %d 'host' must be set", i+1)
		}
		if test.TestType == "" {
			return Config{}, fmt.Errorf("test %d 'testType' must be set", i+1)
		}
	}

	return config, nil
}
