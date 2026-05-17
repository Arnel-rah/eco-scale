package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type ContainerPolicy struct {
	ContainerName string `yaml:"container_name"`
	Priority      string `yaml:"priority"`
	CPULimit      int    `yaml:"cpu_limit"`
}

type Config struct {
	Version   string            `yaml:"version"`
	AlertMode bool              `yaml:"alert_mode"`
	Policies  []ContainerPolicy `yaml:"policies"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
