package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir, configFile), nil
}

func loadConfig() (*GlobalConfig, error) {
	cfgPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	cfg := &GlobalConfig{AppsURL: defaultAppsURL}

	data, err := os.ReadFile(cfgPath)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func saveConfig(cfg *GlobalConfig) error {
	cfgPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	configDir := filepath.Dir(cfgPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(cfgPath, data, 0644)
}
