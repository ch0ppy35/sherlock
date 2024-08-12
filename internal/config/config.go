package config

import (
	"fmt"

	"github.com/ch0ppy35/sherlock/internal/ui"
	"github.com/spf13/viper"
)

type Config struct {
	DNSServer string          `yaml:"dnsServer"`
	Tests     []DNSTestConfig `yaml:"tests"`
}

type DNSTestConfig struct {
	ExpectedValues []string `yaml:"expectedValues"`
	Host           string   `yaml:"host"`
	TestType       string   `yaml:"testType"`
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

	if config.DNSServer == "" {
		ui.PrintMsgWithStatus("WARN", "hiYellow", "DNS server not set, using Cloudflare as default (1.1.1.1)\n")
		config.DNSServer = "1.1.1.1"
	}
	return config, nil
}
