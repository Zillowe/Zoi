package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name     string   `yaml:"name"`
	Guides   []string `yaml:"guides"`
	Provider string   `yaml:"provider"`
	Model    string   `yaml:"model"`
	APIKey   string   `yaml:"api"`
	Endpoint string   `yaml:"endpoint,omitempty"`
}

func LoadConfig() (*Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	for {
		configPath := filepath.Join(dir, "gct.yaml")
		if _, err := os.Stat(configPath); err == nil {
			data, err := os.ReadFile(configPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
			}

			var config Config
			err = yaml.Unmarshal(data, &config)
			if err != nil {
				return nil, fmt.Errorf("failed to parse yaml config %s: %w", configPath, err)
			}
			return &config, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return nil, fmt.Errorf("no gct.yaml file found in this directory or any parent directory")
		}
		dir = parentDir
	}
}
