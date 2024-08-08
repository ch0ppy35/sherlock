package testing

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DNSServer string       `yaml:"dnsServer"`
	Tests     []TestConfig `yaml:"tests"`
}

type TestConfig struct {
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

	return config, nil
}
