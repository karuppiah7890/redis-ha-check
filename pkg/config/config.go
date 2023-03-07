package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Sentinel struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type SentinelConfig struct {
	MasterName string     `yaml:"masterName"`
	Sentinels  []Sentinel `yaml:"sentinels"`
}

type Config struct {
	SentinelConfig SentinelConfig `yaml:"sentinelConfig"`
}

func LoadConfig(configFilePath string) (*Config, error) {
	configFileContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config yaml file at path %s: %v", configFilePath, err)
	}

	var config Config

	err = yaml.Unmarshal(configFileContent, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config yaml file at path %s: %v", configFilePath, err)
	}

	return &config, nil
}
