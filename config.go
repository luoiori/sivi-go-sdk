package sivi

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Sivi SiviConfig `yaml:"sivi"`
}

type SiviConfig struct {
	SDK SDKConfig `yaml:"sdk"`
}

type SDKConfig struct {
	App       string `yaml:"app"`
	AppID     int    `yaml:"app-id"`
	Server    string `yaml:"server"`
	Profile   string `yaml:"profile"`
	MetricURL string `yaml:"metric-url"`
	Period    int    `yaml:"period"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) GetExportInterval() time.Duration {
	return 5 * time.Second
}
