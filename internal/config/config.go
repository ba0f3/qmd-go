package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var CurrentIndexName = "index"

func GetConfigDir() (string, error) {
	if dir := os.Getenv("QMD_CONFIG_DIR"); dir != "" {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "qmd"), nil
}

func GetConfigFilePath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, fmt.Sprintf("%s.yml", CurrentIndexName)), nil
}

func EnsureConfigDir() error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0755)
}

func LoadConfig() (*Config, error) {
	path, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Config{Collections: make(map[string]Collection)}, nil
	}
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Collections == nil {
		cfg.Collections = make(map[string]Collection)
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}
	path, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
